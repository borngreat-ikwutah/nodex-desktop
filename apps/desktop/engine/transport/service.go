package transport

import (
	"context"
	"crypto/ed25519"
	"time"

	"nodex-desktop/engine/identity"
)

// Service provides the Wails facade and runtime lifecycle for Engine 03 (Transport).
type Service struct {
	server *Server
	client *Client
}

// NewService creates a new transport service instance.
func NewService(id identity.Public, signingKey ...ed25519.PrivateKey) *Service {
	var key ed25519.PrivateKey
	if len(signingKey) > 0 {
		key = signingKey[0]
	}
	return &Service{
		server: newServer(id, key),
		client: newClient(id, key),
	}
}

// PingVerified performs a ping while pinning the remote certificate to the
// public identity obtained from discovery.
func (s *Service) PingVerified(ctx context.Context, addr string, expected identity.Public) (*PingResult, error) {
	return s.client.PingVerified(ctx, addr, expected)
}

// Startup starts the TCP transport server.
func (s *Service) Startup(ctx context.Context) error {
	return s.server.Start(ctx)
}

// Shutdown gracefully terminates the TCP transport server.
func (s *Service) Shutdown(ctx context.Context) error {
	return s.server.Stop()
}

// Ping connects to the remote address (e.g. "192.168.1.50:4242"),
// sends a framed ping, and measures round-trip latency.
func (s *Service) Ping(addr string) (*PingResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Ping(ctx, addr)
}

// Port returns the active TCP listening port.
func (s *Service) Port() int {
	return s.server.Port()
}

// ActiveConnections returns the number of active peer TCP connections.
func (s *Service) ActiveConnections() int {
	return s.server.ActiveConnections()
}

// Server returns the underlying TCP server.
func (s *Service) Server() *Server {
	return s.server
}

// Client returns the underlying TCP client.
func (s *Service) Client() *Client {
	return s.client
}
