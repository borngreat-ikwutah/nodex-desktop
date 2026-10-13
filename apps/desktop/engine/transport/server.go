package transport

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"nodex-desktop/engine/identity"
)

const (
	DefaultServerPort = 4242
)

// HandlerFunc processes incoming TCP messages and optionally returns a response message.
type HandlerFunc func(conn net.Conn, msg *Message) (*Message, error)

// ServerOption configures Server parameters.
type ServerOption func(*Server)

// WithServerPort overrides the default TCP listening port.
func WithServerPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

// WithServerListenIP sets the IP to bind the TCP listener to (default "0.0.0.0").
func WithServerListenIP(ip string) ServerOption {
	return func(s *Server) {
		s.listenIP = ip
	}
}

// Server is the TCP listening engine for Nodex peer-to-peer connections.
type Server struct {
	identity   identity.Public
	signingKey ed25519.PrivateKey
	port       int
	listenIP   string

	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	connsMu sync.Mutex
	conns   map[net.Conn]struct{}

	handlersMu sync.RWMutex
	handlers   map[MessageType]HandlerFunc
}

// NewServer initializes a TCP server backed by this machine's public identity.
func NewServer(id identity.Public, opts ...ServerOption) *Server {
	return newServer(id, nil, opts...)
}

func newServer(id identity.Public, signingKey ed25519.PrivateKey, opts ...ServerOption) *Server {
	s := &Server{
		identity:   id,
		signingKey: signingKey,
		port:       DefaultServerPort,
		listenIP:   "0.0.0.0",
		conns:      make(map[net.Conn]struct{}),
		handlers:   make(map[MessageType]HandlerFunc),
	}
	for _, opt := range opts {
		opt(s)
	}

	// Register default ping handler
	s.RegisterHandler(TypePing, func(conn net.Conn, msg *Message) (*Message, error) {
		pong := NewMessage(TypePong, s.identity.NodeID, s.identity.UUID, s.identity.PublicKey)
		pong.TargetID = msg.SenderID
		pong.Payload = msg.Payload
		return pong, nil
	})

	return s
}

// RegisterHandler registers a custom message handler for a specific message type.
func (s *Server) RegisterHandler(msgType MessageType, handler HandlerFunc) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.handlers[msgType] = handler
}

// Start begins listening for incoming TCP connections.
func (s *Server) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	addr := fmt.Sprintf("%s:%d", s.listenIP, s.port)
	var listener net.Listener
	if len(s.signingKey) == ed25519.PrivateKeySize {
		cert, err := newIdentityCertificate(s.identity, s.signingKey)
		if err != nil {
			return err
		}
		listener, err = tls.Listen("tcp4", addr, serverTLSConfig(cert))
		if err != nil {
			return fmt.Errorf("transport: failed to bind TLS listener on %s: %w", addr, err)
		}
	} else {
		var err error
		listener, err = net.Listen("tcp4", addr)
		if err != nil {
			return fmt.Errorf("transport: failed to bind TCP listener on %s: %w", addr, err)
		}
	}
	s.listener = listener

	// Update actual port if bound to dynamic port 0
	if tcpAddr, ok := listener.Addr().(*net.TCPAddr); ok {
		s.port = tcpAddr.Port
	}

	s.wg.Add(1)
	go s.acceptLoop()

	return nil
}

// Port returns the active TCP listening port.
func (s *Server) Port() int {
	return s.port
}

// ActiveConnections returns the number of currently connected peer sockets.
func (s *Server) ActiveConnections() int {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()
	return len(s.conns)
}

// acceptLoop accepts new TCP connections until the server is closed.
func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				log.Printf("transport: TCP accept error: %v\n", err)
				continue
			}
		}

		s.trackConn(conn)
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			defer s.untrackConn(c)
			s.handleConn(c)
		}(conn)
	}
}

// handleConn reads framed messages in a loop and executes registered handlers.
func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		msg, err := ReadMessage(conn)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			select {
			case <-s.ctx.Done():
				return
			default:
				log.Printf("transport: connection read error from %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}

		s.handlersMu.RLock()
		handler, exists := s.handlers[msg.Type]
		s.handlersMu.RUnlock()

		if !exists {
			log.Printf("transport: unhandled message type %q from %s\n", msg.Type, conn.RemoteAddr())
			continue
		}

		reply, err := handler(conn, msg)
		if err != nil {
			log.Printf("transport: handler error for %s: %v\n", msg.Type, err)
			continue
		}

		if reply != nil {
			if err := WriteMessage(conn, reply); err != nil {
				log.Printf("transport: failed to write reply to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}
	}
}

func (s *Server) trackConn(c net.Conn) {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()
	s.conns[c] = struct{}{}
}

func (s *Server) untrackConn(c net.Conn) {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()
	delete(s.conns, c)
}

// Stop shuts down the TCP listener and closes all active peer connections.
func (s *Server) Stop() error {
	if s.cancel != nil {
		s.cancel()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	s.connsMu.Lock()
	for c := range s.conns {
		_ = c.Close()
	}
	s.connsMu.Unlock()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
	}

	return nil
}
