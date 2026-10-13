package transport

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"nodex-desktop/engine/identity"
)

const (
	DefaultDialTimeout = 3 * time.Second
)

// PingResult contains metrics and identity details returned from a successful ping.
type PingResult struct {
	PeerID        string        `json:"peer_id"`
	PeerUUID      string        `json:"peer_uuid"`
	PeerPublicKey string        `json:"peer_public_key"`
	Latency       time.Duration `json:"latency"`
	LatencyMs     float64       `json:"latency_ms"`
	RemoteAddr    string        `json:"remote_addr"`
}

// Client establishes outbound TCP transport connections to remote Nodex nodes.
type Client struct {
	identity   identity.Public
	signingKey ed25519.PrivateKey
	timeout    time.Duration
}

// NewClient initializes a transport client with the local machine identity.
func NewClient(id identity.Public, timeout ...time.Duration) *Client {
	return newClient(id, nil, timeout...)
}

func newClient(id identity.Public, signingKey ed25519.PrivateKey, timeout ...time.Duration) *Client {
	t := DefaultDialTimeout
	if len(timeout) > 0 && timeout[0] > 0 {
		t = timeout[0]
	}
	return &Client{
		identity:   id,
		signingKey: signingKey,
		timeout:    t,
	}
}

// Dial connects to a remote peer's TCP transport socket.
func (c *Client) Dial(ctx context.Context, targetAddr string) (net.Conn, error) {
	return c.dial(ctx, targetAddr, nil)
}

func (c *Client) dial(ctx context.Context, targetAddr string, expected *identity.Public) (net.Conn, error) {
	var d net.Dialer
	d.Timeout = c.timeout
	if len(c.signingKey) != ed25519.PrivateKeySize {
		conn, err := d.DialContext(ctx, "tcp4", targetAddr)
		if err != nil {
			return nil, fmt.Errorf("transport: failed to dial peer %s: %w", targetAddr, err)
		}
		return conn, nil
	}
	cert, err := newIdentityCertificate(c.identity, c.signingKey)
	if err != nil {
		return nil, err
	}
	conn, err := tls.DialWithDialer(&d, "tcp4", targetAddr, clientTLSConfig(cert, expected))
	if err != nil {
		return nil, fmt.Errorf("transport: failed to dial peer %s: %w", targetAddr, err)
	}
	return conn, nil
}

// Ping sends a length-prefixed ping message to targetAddr and awaits a pong response,
// measuring round-trip network latency.
func (c *Client) Ping(ctx context.Context, targetAddr string) (*PingResult, error) {
	return c.ping(ctx, targetAddr, nil)
}

// PingVerified pins the peer certificate to the identity learned during discovery.
func (c *Client) PingVerified(ctx context.Context, targetAddr string, expected identity.Public) (*PingResult, error) {
	return c.ping(ctx, targetAddr, &expected)
}

func (c *Client) ping(ctx context.Context, targetAddr string, expected *identity.Public) (*PingResult, error) {
	conn, err := c.dial(ctx, targetAddr, expected)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Apply deadline from context if present
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(c.timeout))
	}

	pingMsg := NewMessage(TypePing, c.identity.NodeID, c.identity.UUID, c.identity.PublicKey)
	start := time.Now()

	if err := WriteMessage(conn, pingMsg); err != nil {
		return nil, fmt.Errorf("transport: failed to send ping frame: %w", err)
	}

	reply, err := ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("transport: failed to receive pong frame: %w", err)
	}

	latency := time.Since(start)

	if reply.Type != TypePong {
		return nil, fmt.Errorf("transport: expected pong message, received %q", reply.Type)
	}

	return &PingResult{
		PeerID:        reply.SenderID,
		PeerUUID:      reply.UUID,
		PeerPublicKey: reply.PublicKey,
		Latency:       latency,
		LatencyMs:     float64(latency.Microseconds()) / 1000.0,
		RemoteAddr:    targetAddr,
	}, nil
}
