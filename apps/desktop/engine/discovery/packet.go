package discovery

import (
	"encoding/json"
	"errors"
	"strings"
	"unicode"
)

var (
	// ErrNodeNotFound is returned when a node does not reply within the discovery deadline.
	ErrNodeNotFound = errors.New("node not found on local network")
	// ErrInvalidNodeID is returned when the provided node ID is malformed.
	ErrInvalidNodeID = errors.New("invalid 9-digit node id")
)

// PacketType defines the role of a discovery datagram.
type PacketType string

const (
	PacketDiscover PacketType = "discover" // Shout: "Is Node X here?"
	PacketReply    PacketType = "reply"    // Reply: "Yes, I am Node X!"
)

// Packet represents the wire format of UDP discovery messages.
type Packet struct {
	Type      PacketType `json:"type"`                 // "discover" or "reply"
	TargetID  string     `json:"target_id"`            // Who we are looking for (e.g. "222 222 222")
	SenderID  string     `json:"sender_id"`            // Who is asking / replying (e.g. "111 111 111")
	SenderIP  string     `json:"sender_ip,omitempty"`  // Reported IP of the sender
	UUID      string     `json:"uuid,omitempty"`       // Unique machine UUID (from Engine 01)
	PublicKey string     `json:"public_key,omitempty"` // ed25519 public key hex for transport verification
	Port      int        `json:"port,omitempty"`       // Service/Transport port (e.g. 4242)
}

// Marshal encodes a packet into JSON bytes.
func (p *Packet) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalPacket decodes JSON bytes into a Packet.
func UnmarshalPacket(data []byte) (*Packet, error) {
	var p Packet
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// NormalizeNodeID strips spaces and non-digit characters, formatting it as "XXX XXX XXX".
// Returns an error if the ID does not contain exactly 9 digits.
func NormalizeNodeID(raw string) (string, error) {
	var digits strings.Builder
	for _, r := range raw {
		if unicode.IsDigit(r) {
			digits.WriteRune(r)
		}
	}
	s := digits.String()
	if len(s) != 9 {
		return "", ErrInvalidNodeID
	}
	return s[:3] + " " + s[3:6] + " " + s[6:], nil
}
