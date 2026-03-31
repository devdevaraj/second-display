package stream

import (
	"fmt"
	"log"
	"net"

	"github.com/pion/webrtc/v4"
)

// WebRTCHandler processes incoming RTP stream from FFmpeg and serves WebRTC clients.
type WebRTCHandler struct {
	SessionID  string
	RTPPort    int
	VideoTrack *webrtc.TrackLocalStaticRTP
	conn       *net.UDPConn
}

// NewWebRTCHandler initializes the track and WebRTC handler.
func NewWebRTCHandler(sessionID string, rtpPort int) (*WebRTCHandler, error) {
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
	if err != nil {
		return nil, err
	}

	h := &WebRTCHandler{
		SessionID:  sessionID,
		RTPPort:    rtpPort,
		VideoTrack: videoTrack,
	}
	return h, nil
}

// StartRTPListener opens a local UDP port to receive RTP packets from FFmpeg.
func (h *WebRTCHandler) StartRTPListener() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", h.RTPPort))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	h.conn = conn

	go func() {
		defer conn.Close()
		buf := make([]byte, 1500)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				return // Closed
			}
			if _, writeErr := h.VideoTrack.Write(buf[:n]); writeErr != nil {
				log.Printf("Error writing RTP packet: %v", writeErr)
			}
		}
	}()
	return nil
}

// Stop closes the RTP listener socket.
func (h *WebRTCHandler) Stop() {
	if h.conn != nil {
		h.conn.Close()
	}
}

// HandleOffer processes an incoming SDP Offer and generates an SDP Answer.
func (h *WebRTCHandler) HandleOffer(offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	api := webrtc.NewAPI()
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, err
	}

	if _, err = peerConnection.AddTrack(h.VideoTrack); err != nil {
		return nil, err
	}

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		return nil, err
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	// Wait for ICE Gathering to complete before returning answer
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	if err = peerConnection.SetLocalDescription(answer); err != nil {
		return nil, err
	}

	<-gatherComplete

	return peerConnection.LocalDescription(), nil
}
