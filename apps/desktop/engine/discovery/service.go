package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"nodex-desktop/engine/identity"
)

const (
	DefaultDiscoveryPort = 9999
	DefaultTransportPort = 4242
	DefaultRetryInterval = 350 * time.Millisecond
	DefaultProbeTimeout  = 4 * time.Second
)

// Peer represents another Nodex installation discovered on the LAN.
type Peer struct {
	NodeID    string   `json:"node_id"`
	UUID      string   `json:"uuid"`
	PublicKey string   `json:"public_key"`
	IPs       []string `json:"ips"`
	Port      int      `json:"port"`
}

// Option configures discovery Service parameters (useful for tests or custom topologies).
type Option func(*Service)

// WithPort overrides the UDP broadcast and listening port.
func WithPort(port int) Option {
	return func(s *Service) {
		s.port = port
	}
}

// WithListenIP overrides the IP to bind the UDP listener to (e.g. "127.0.0.1" for tests).
func WithListenIP(ip string) Option {
	return func(s *Service) {
		s.listenIP = ip
	}
}

// WithExtraBroadcastTargets adds explicit UDP addresses to broadcast to (e.g. "127.0.0.1:9999" in tests).
func WithExtraBroadcastTargets(targets ...string) Option {
	return func(s *Service) {
		s.extraTargets = append(s.extraTargets, targets...)
	}
}

// WithTransportPort sets the listening port of the transport engine.
func WithTransportPort(port int) Option {
	return func(s *Service) {
		s.transportPort = port
	}
}

// Service provides LAN discovery of other Nodex instances over UDP broadcast.
type Service struct {
	ctx           context.Context
	cancel        context.CancelFunc
	identity      identity.Public
	port          int
	transportPort int
	listenIP      string
	extraTargets  []string

	listenerConn *net.UDPConn
	listenerDone chan struct{}

	peersMu sync.RWMutex
	peers   map[string]Peer // Keyed by UUID

	waitersMu sync.Mutex
	waiters   map[string][]chan Peer // Keyed by normalized Target Node ID
}

// NewService creates a new UDP LAN discovery engine.
func NewService(id identity.Public, opts ...Option) *Service {
	s := &Service{
		identity:      id,
		port:          DefaultDiscoveryPort,
		transportPort: DefaultTransportPort,
		listenIP:      "0.0.0.0",
		peers:         make(map[string]Peer),
		waiters:       make(map[string][]chan Peer),
		listenerDone:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Startup starts the UDP background listener ("The Ear").
func (s *Service) Startup(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	listenAddr := &net.UDPAddr{
		IP:   net.ParseIP(s.listenIP),
		Port: s.port,
	}

	conn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		return fmt.Errorf("discovery: failed to bind UDP listener on %s:%d: %w", s.listenIP, s.port, err)
	}
	s.listenerConn = conn

	go s.listenLoop()

	return nil
}

// listenLoop continuously receives UDP datagrams, answering shouts and updating peers.
func (s *Service) listenLoop() {
	defer close(s.listenerDone)

	buf := make([]byte, 4096)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		n, remoteAddr, err := s.listenerConn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				log.Printf("discovery: UDP read error: %v\n", err)
				continue
			}
		}

		packet, err := UnmarshalPacket(buf[:n])
		if err != nil {
			// Malformed or foreign UDP packet, safely discard
			continue
		}

		// Self-reflection check: ignore packets shouted or replied by our own node
		if packet.SenderID == s.identity.NodeID || (packet.UUID != "" && packet.UUID == s.identity.UUID) {
			continue
		}

		switch packet.Type {
		case PacketDiscover:
			s.handleDiscover(packet, remoteAddr)
		case PacketReply:
			s.handleReply(packet, remoteAddr)
		}
	}
}

// handleDiscover responds to a peer asking if our Node ID is on this network.
func (s *Service) handleDiscover(packet *Packet, remoteAddr *net.UDPAddr) {
	normTarget, err := NormalizeNodeID(packet.TargetID)
	if err != nil {
		return
	}

	// Only reply if the shout is asking for our Node ID
	if normTarget != s.identity.NodeID {
		return
	}

	outboundIP := getPreferredOutboundIP()
	reply := Packet{
		Type:      PacketReply,
		TargetID:  packet.SenderID,
		SenderID:  s.identity.NodeID,
		SenderIP:  outboundIP,
		UUID:      s.identity.UUID,
		PublicKey: s.identity.PublicKey,
		Port:      s.transportPort,
	}

	data, err := reply.Marshal()
	if err != nil {
		return
	}

	// Whisper back directly to the sender's remote UDP address
	if _, err := s.listenerConn.WriteToUDP(data, remoteAddr); err != nil {
		log.Printf("discovery: failed to reply to %s: %v\n", remoteAddr, err)
	}
}

// handleReply processes incoming reply packets from target nodes.
func (s *Service) handleReply(packet *Packet, remoteAddr *net.UDPAddr) {
	normTarget, err := NormalizeNodeID(packet.TargetID)
	if err != nil {
		return
	}

	// Verify the reply was addressed to us
	if normTarget != s.identity.NodeID {
		return
	}

	resolvedIP := remoteAddr.IP.String()
	if packet.SenderIP != "" {
		resolvedIP = packet.SenderIP
	}

	peer := Peer{
		NodeID:    packet.SenderID,
		UUID:      packet.UUID,
		PublicKey: packet.PublicKey,
		IPs:       []string{resolvedIP},
		Port:      packet.Port,
	}

	// Cache peer
	peerKey := peer.UUID
	if peerKey == "" {
		peerKey = peer.NodeID
	}

	s.peersMu.Lock()
	s.peers[peerKey] = peer
	s.peersMu.Unlock()

	// Notify active FindNode waiters
	normSender, _ := NormalizeNodeID(packet.SenderID)
	s.waitersMu.Lock()
	if chList, ok := s.waiters[normSender]; ok {
		for _, ch := range chList {
			select {
			case ch <- peer:
			default:
			}
		}
	}
	s.waitersMu.Unlock()

	// Notify Wails frontend
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "peer:discovered", peer)
	}
}

// FindNode locates the target machine's local IP address on the local network given its 9-digit Node ID.
func (s *Service) FindNode(ctx context.Context, targetID string) (*Peer, error) {
	normTarget, err := NormalizeNodeID(targetID)
	if err != nil {
		return nil, err
	}

	if normTarget == s.identity.NodeID {
		return nil, fmt.Errorf("cannot discover self (%s)", normTarget)
	}

	// Register waiter channel
	replyChan := make(chan Peer, 4)
	s.waitersMu.Lock()
	s.waiters[normTarget] = append(s.waiters[normTarget], replyChan)
	s.waitersMu.Unlock()

	defer func() {
		s.waitersMu.Lock()
		waiterList := s.waiters[normTarget]
		updated := make([]chan Peer, 0, len(waiterList))
		for _, ch := range waiterList {
			if ch != replyChan {
				updated = append(updated, ch)
			}
		}
		if len(updated) == 0 {
			delete(s.waiters, normTarget)
		} else {
			s.waiters[normTarget] = updated
		}
		s.waitersMu.Unlock()
	}()

	// Open ephemeral UDP socket for broadcasting & reading direct replies
	bcastConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, fmt.Errorf("discovery: failed to open broadcast socket: %w", err)
	}
	defer bcastConn.Close()

	packet := Packet{
		Type:      PacketDiscover,
		TargetID:  normTarget,
		SenderID:  s.identity.NodeID,
		SenderIP:  getPreferredOutboundIP(),
		UUID:      s.identity.UUID,
		PublicKey: s.identity.PublicKey,
		Port:      s.transportPort,
	}

	payload, err := packet.Marshal()
	if err != nil {
		return nil, err
	}

	// Prepare list of broadcast destination UDP addresses
	destAddrs := make([]*net.UDPAddr, 0)
	for _, ip := range getBroadcastAddresses() {
		destAddrs = append(destAddrs, &net.UDPAddr{IP: ip, Port: s.port})
	}
	for _, target := range s.extraTargets {
		if udpAddr, err := net.ResolveUDPAddr("udp4", target); err == nil {
			destAddrs = append(destAddrs, udpAddr)
		}
	}

	sendBroadcast := func() {
		for _, dest := range destAddrs {
			_, _ = bcastConn.WriteToUDP(payload, dest)
		}
	}

	// Background worker reading replies on ephemeral socket
	go func() {
		readBuf := make([]byte, 2048)
		for {
			_ = bcastConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			n, rAddr, err := bcastConn.ReadFromUDP(readBuf)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}

			reply, err := UnmarshalPacket(readBuf[:n])
			if err != nil || reply.Type != PacketReply {
				continue
			}

			replySender, _ := NormalizeNodeID(reply.SenderID)
			if replySender == normTarget {
				resolvedIP := rAddr.IP.String()
				if reply.SenderIP != "" {
					resolvedIP = reply.SenderIP
				}
				peer := Peer{
					NodeID:    reply.SenderID,
					UUID:      reply.UUID,
					PublicKey: reply.PublicKey,
					IPs:       []string{resolvedIP},
					Port:      reply.Port,
				}
				select {
				case replyChan <- peer:
				default:
				}
				return
			}
		}
	}()

	// Send initial shout
	sendBroadcast()

	ticker := time.NewTicker(DefaultRetryInterval)
	defer ticker.Stop()

	for {
		select {
		case peer := <-replyChan:
			return &peer, nil
		case <-ticker.C:
			// Retransmit shout in case UDP packet was dropped on Wi-Fi
			sendBroadcast()
		case <-ctx.Done():
			return nil, ErrNodeNotFound
		}
	}
}

// Shutdown stops the UDP background listener and closes network sockets.
func (s *Service) Shutdown(ctx context.Context) {
	if s.cancel != nil {
		s.cancel()
	}
	if s.listenerConn != nil {
		_ = s.listenerConn.Close()
	}

	select {
	case <-s.listenerDone:
	case <-time.After(500 * time.Millisecond):
	}
}

// GetPeers returns all currently discovered peers on the LAN.
func (s *Service) GetPeers() []Peer {
	s.peersMu.RLock()
	defer s.peersMu.RUnlock()

	res := make([]Peer, 0, len(s.peers))
	for _, p := range s.peers {
		res = append(res, p)
	}
	return res
}
