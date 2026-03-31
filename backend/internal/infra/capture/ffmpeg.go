package capture

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
)

// Pipeline manages the wf-recorder and ffmpeg processes for a session.
type Pipeline struct {
	SessionID  string
	OutputName string
	Resolution string
	FPS        int
	Bitrate    int
	RTPPort    int // Port for WebRTC backend to consume

	cmdRecorder *exec.Cmd
	cmdFFmpeg   *exec.Cmd
	cancel      context.CancelFunc
	mu          sync.Mutex
}

// NewPipeline initializes a new capture pipeline.
func NewPipeline(sessionID, outputName, resolution string, fps, bitrate, rtpPort int) *Pipeline {
	return &Pipeline{
		SessionID:  sessionID,
		OutputName: outputName,
		Resolution: resolution,
		FPS:        fps,
		Bitrate:    bitrate,
		RTPPort:    rtpPort,
	}
}

// Start launches the capture processes. wf-recorder writes to stdout, mapped to ffmpeg native RTP.
func (p *Pipeline) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	// Start wf-recorder targeting the specific headless output and outputting rawvideo.
	// We pipe this to ffmpeg.
	p.cmdRecorder = exec.CommandContext(ctx, "wf-recorder",
		"-o", p.OutputName,
		"-c", "rawvideo",
		"-m", "avi", // Using AVI to stream to stdout without needing seeking
		"-f", "pipe:1",
	)

	// FFmpeg reads from stdin and outputs an RTP stream for Pion WebRTC to consume locally.
	// We use h264 software encoder (libx264) with ultrafast preset for low latency.
	bitrateStr := fmt.Sprintf("%dk", p.Bitrate)
	rtpDest := fmt.Sprintf("rtp://127.0.0.1:%d", p.RTPPort)

	p.cmdFFmpeg = exec.CommandContext(ctx, "ffmpeg",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-b:v", bitrateStr,
		"-r", fmt.Sprintf("%d", p.FPS),
		"-f", "rtp",
		rtpDest,
	)

	// Pipe recorder stdout to ffmpeg stdin
	pipeReader, pipeWriter := io.Pipe()
	p.cmdRecorder.Stdout = pipeWriter
	p.cmdFFmpeg.Stdin = pipeReader

	// Capture errors for logging
	p.cmdRecorder.Stderr = log.Writer()
	p.cmdFFmpeg.Stderr = log.Writer()

	// Start ffmpeg first so it's ready to read
	if err := p.cmdFFmpeg.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Start wf-recorder
	if err := p.cmdRecorder.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start wf-recorder: %w", err)
	}

	// Wait asynchronously to handle teardown
	go func() {
		_ = p.cmdRecorder.Wait()
		pipeWriter.Close() // notify ffmpeg of EOF
		_ = p.cmdFFmpeg.Wait()
	}()

	return nil
}

// Stop gracefully stops the pipeline avoiding zombies.
func (p *Pipeline) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
}
