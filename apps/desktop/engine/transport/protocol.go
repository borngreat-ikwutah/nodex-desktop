package transport

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	// MaxPayloadSize sets an upper limit on frame size (16 MB) to prevent OOM attacks.
	MaxPayloadSize = 16 * 1024 * 1024

	// HeaderSize is the 4-byte uint32 length prefix.
	HeaderSize = 4
)

var (
	// ErrMessageTooLarge is returned when an incoming or outgoing frame exceeds MaxPayloadSize.
	ErrMessageTooLarge = errors.New("transport: frame payload exceeds maximum allowed size (16MB)")
	// ErrInvalidHeader is returned when reading the frame header fails or is incomplete.
	ErrInvalidHeader = errors.New("transport: failed to read complete 4-byte frame header")
)

// MessageType distinguishes control and data messages in Engine 03.
type MessageType string

const (
	TypePing         MessageType = "ping"
	TypePong         MessageType = "pong"
	TypeData         MessageType = "data"
	TypeWebRTCOffer  MessageType = "webrtc_offer"
	TypeWebRTCAnswer MessageType = "webrtc_answer"
	TypeWebRTCICE    MessageType = "webrtc_ice"
	TypeWebRTCClose  MessageType = "webrtc_close"
)

// Message is the standard framed envelope exchanged over Nodex TCP sockets.
type Message struct {
	Type      MessageType     `json:"type"`                 // "ping" | "pong" | "data"
	SenderID  string          `json:"sender_id,omitempty"`  // Sender Node ID
	TargetID  string          `json:"target_id,omitempty"`  // Target Node ID
	UUID      string          `json:"uuid,omitempty"`       // Machine UUID from Engine 01
	PublicKey string          `json:"public_key,omitempty"` // ed25519 public key
	Timestamp int64           `json:"timestamp"`            // Unix nano timestamp
	Payload   json.RawMessage `json:"payload,omitempty"`    // Optional arbitrary data
}

// NewMessage creates a message with the current timestamp.
func NewMessage(msgType MessageType, senderID, uuid, pubKey string) *Message {
	return &Message{
		Type:      msgType,
		SenderID:  senderID,
		UUID:      uuid,
		PublicKey: pubKey,
		Timestamp: time.Now().UnixNano(),
	}
}

// WriteMessage marshals msg to JSON, prepends a 4-byte big-endian length prefix,
// and writes the complete frame to w.
func WriteMessage(w io.Writer, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("transport: failed to marshal message: %w", err)
	}

	payloadLen := len(data)
	if payloadLen > MaxPayloadSize {
		return ErrMessageTooLarge
	}

	var header [HeaderSize]byte
	binary.BigEndian.PutUint32(header[:], uint32(payloadLen))

	// Write length prefix
	if _, err := w.Write(header[:]); err != nil {
		return fmt.Errorf("transport: failed to write frame header: %w", err)
	}

	// Write payload
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("transport: failed to write frame payload: %w", err)
	}

	return nil
}

// ReadMessage reads a 4-byte length header from r using io.ReadFull,
// reads exactly that many payload bytes, and unmarshals into a Message.
func ReadMessage(r io.Reader) (*Message, error) {
	var header [HeaderSize]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidHeader, err)
	}

	payloadLen := binary.BigEndian.Uint32(header[:])
	if payloadLen > MaxPayloadSize {
		return nil, ErrMessageTooLarge
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf("transport: failed to read complete frame payload: %w", err)
	}

	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, fmt.Errorf("transport: failed to unmarshal message payload: %w", err)
	}

	return &msg, nil
}
