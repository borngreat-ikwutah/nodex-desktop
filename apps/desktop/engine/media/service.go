package media

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pion/webrtc/v4"

	"nodex-desktop/engine/identity"
)

// Signal is the signaling payload exchanged over the authenticated transport.
// The media packets themselves do not pass through Engine 3.
type Signal struct {
	SessionID string  `json:"session_id"`
	SDP       string  `json:"sdp,omitempty"`
	Candidate string  `json:"candidate,omitempty"`
	SDPMLine  *uint16 `json:"sdp_mline_index,omitempty"`
}

// Service owns WebRTC sessions for this installation.
type Service struct {
	identity identity.Public
	api      *webrtc.API
	mu       sync.Mutex
	sessions map[string]*Session
}

// Session is a single authenticated WebRTC media negotiation.
type Session struct {
	ID string
	PC *webrtc.PeerConnection
}

func NewService(id identity.Public) *Service {
	mediaEngine := &webrtc.MediaEngine{}
	_ = mediaEngine.RegisterDefaultCodecs()
	return &Service{
		identity: id,
		api:      webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine)),
		sessions: make(map[string]*Session),
	}
}

// CreateOffer creates a host peer connection and returns its SDP offer.
func (s *Service) CreateOffer(ctx context.Context, sessionID string) (Signal, error) {
	if sessionID == "" {
		return Signal{}, fmt.Errorf("media: session id is required")
	}
	pc, err := s.api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return Signal{}, fmt.Errorf("media: create peer connection: %w", err)
	}
	if _, err := pc.CreateDataChannel("nodex-control", nil); err != nil {
		_ = pc.Close()
		return Signal{}, fmt.Errorf("media: create control channel: %w", err)
	}
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		_ = pc.Close()
		return Signal{}, fmt.Errorf("media: create offer: %w", err)
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		_ = pc.Close()
		return Signal{}, fmt.Errorf("media: set local offer: %w", err)
	}
	s.mu.Lock()
	s.sessions[sessionID] = &Session{ID: sessionID, PC: pc}
	s.mu.Unlock()
	return Signal{SessionID: sessionID, SDP: offer.SDP}, nil
}

// AcceptAnswer completes the host side negotiation.
func (s *Service) AcceptAnswer(_ context.Context, signal Signal) error {
	session, err := s.session(signal.SessionID)
	if err != nil {
		return err
	}
	return session.PC.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: signal.SDP})
}

// AddICECandidate adds a trickled candidate to an existing session.
func (s *Service) AddICECandidate(_ context.Context, signal Signal) error {
	session, err := s.session(signal.SessionID)
	if err != nil {
		return err
	}
	return session.PC.AddICECandidate(webrtc.ICECandidateInit{Candidate: signal.Candidate, SDPMLineIndex: signal.SDPMLine})
}

func (s *Service) Close(sessionID string) error {
	s.mu.Lock()
	session, ok := s.sessions[sessionID]
	if ok {
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()
	if !ok {
		return nil
	}
	return session.PC.Close()
}

func (s *Service) session(id string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("media: unknown session %q", id)
	}
	return session, nil
}

func EncodeSignal(signal Signal) (json.RawMessage, error) { return json.Marshal(signal) }

func DecodeSignal(payload json.RawMessage) (Signal, error) {
	var signal Signal
	if err := json.Unmarshal(payload, &signal); err != nil {
		return Signal{}, fmt.Errorf("media: decode signal: %w", err)
	}
	return signal, nil
}
