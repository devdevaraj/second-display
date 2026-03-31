package service

import (
	"context"
	"fmt"
	"sync"

	"vdisplay/internal/domain"
	"vdisplay/internal/infra/capture"
	"vdisplay/internal/infra/sway"
	"vdisplay/internal/stream"

	"github.com/google/uuid"
)

type RunningSession struct {
	Details  *domain.Session
	Pipeline *capture.Pipeline
	WebRTC   *stream.WebRTCHandler
	Ctx      context.Context
	Cancel   context.CancelFunc
}

type SessionManager struct {
	mu           sync.RWMutex
	sessions     map[string]*RunningSession
	nextPort     int
	portMu       sync.Mutex
	outputCounts int // for HEADLESS-X
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*RunningSession),
		nextPort: 5000,
	}
}

func (sm *SessionManager) getNextPort() int {
	sm.portMu.Lock()
	defer sm.portMu.Unlock()
	port := sm.nextPort
	sm.nextPort += 2 // RTP/RTCP pair spacing convention
	return port
}

func (sm *SessionManager) getNextOutputName() string {
	sm.portMu.Lock()
	defer sm.portMu.Unlock()
	sm.outputCounts++
	return fmt.Sprintf("HEADLESS-%d", sm.outputCounts)
}

func (sm *SessionManager) CreateSession(resolution string, fps, bitrate int) (*domain.Session, error) {
	id := uuid.New().String()
	outputName := sm.getNextOutputName()

	sess := &domain.Session{
		ID:         id,
		OutputName: outputName,
		Resolution: resolution,
		FPS:        fps,
		Bitrate:    bitrate,
		Status:     domain.StatusCreated,
	}

	sm.mu.Lock()
	sm.sessions[id] = &RunningSession{
		Details: sess,
	}
	sm.mu.Unlock()

	return sess, nil
}

func (sm *SessionManager) StartSession(ctx context.Context, id string) error {
	sm.mu.Lock()
	rs, ok := sm.sessions[id]
	if !ok {
		sm.mu.Unlock()
		return fmt.Errorf("session not found")
	}
	if rs.Details.Status == domain.StatusRunning {
		sm.mu.Unlock()
		return nil
	}
	sessCtx, cancel := context.WithCancel(context.Background())
	rs.Ctx = sessCtx
	rs.Cancel = cancel
	sm.mu.Unlock()

	// 1. Sway Output
	if err := sway.CreateOutput(sessCtx, rs.Details.OutputName, rs.Details.Resolution); err != nil {
		rs.Cancel()
		return err
	}

	// 2. WebRTC Handler
	rtpPort := sm.getNextPort()
	webrtcHandler, err := stream.NewWebRTCHandler(id, rtpPort)
	if err != nil {
		rs.Cancel()
		sway.DestroyOutput(context.Background(), rs.Details.OutputName)
		return err
	}
	if err := webrtcHandler.StartRTPListener(); err != nil {
		rs.Cancel()
		sway.DestroyOutput(context.Background(), rs.Details.OutputName)
		return err
	}
	rs.WebRTC = webrtcHandler

	// 3. Capture Pipeline
	pipeline := capture.NewPipeline(id, rs.Details.OutputName, rs.Details.Resolution, rs.Details.FPS, rs.Details.Bitrate, rtpPort)
	if err := pipeline.Start(sessCtx); err != nil {
		webrtcHandler.Stop()
		sway.DestroyOutput(context.Background(), rs.Details.OutputName)
		rs.Cancel()
		return err
	}
	rs.Pipeline = pipeline

	sm.mu.Lock()
	rs.Details.Status = domain.StatusRunning
	sm.mu.Unlock()

	return nil
}

func (sm *SessionManager) StopSession(id string) error {
	sm.mu.Lock()
	rs, ok := sm.sessions[id]
	if !ok {
		sm.mu.Unlock()
		return fmt.Errorf("session not found")
	}
	if rs.Details.Status != domain.StatusRunning {
		sm.mu.Unlock()
		return nil
	}
	sm.mu.Unlock()

	// Cleanup
	if rs.Cancel != nil {
		rs.Cancel()
	}
	if rs.Pipeline != nil {
		rs.Pipeline.Stop()
	}
	if rs.WebRTC != nil {
		rs.WebRTC.Stop()
	}
	_ = sway.DestroyOutput(context.Background(), rs.Details.OutputName)

	sm.mu.Lock()
	rs.Details.Status = domain.StatusStopped
	sm.mu.Unlock()
	return nil
}

func (sm *SessionManager) DeleteSession(id string) error {
	_ = sm.StopSession(id)

	sm.mu.Lock()
	delete(sm.sessions, id)
	sm.mu.Unlock()
	return nil
}

func (sm *SessionManager) ListSessions() []*domain.Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	list := make([]*domain.Session, 0, len(sm.sessions))
	for _, rs := range sm.sessions {
		list = append(list, rs.Details)
	}
	return list
}

func (sm *SessionManager) GetWebRTC(id string) (*stream.WebRTCHandler, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	rs, ok := sm.sessions[id]
	if !ok || rs.WebRTC == nil {
		return nil, fmt.Errorf("webrtc not available for session")
	}
	return rs.WebRTC, nil
}
