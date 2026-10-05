package discovery

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"nodex-desktop/engine/identity"
)

// Peer represents another Nodex installation discovered on the LAN.
type Peer struct {
	NodeID    string   `json:"node_id"`
	UUID      string   `json:"uuid"`
	PublicKey string   `json:"public_key"`
	IPs       []string `json:"ips"`
	Port      int      `json:"port"`
}

// Service provides LAN discovery of other Nodex instances using mDNS.
type Service struct {
	ctx      context.Context
	identity identity.Public
	server   *zeroconf.Server

	peersMu sync.RWMutex
	peers   map[string]Peer
}

// NewService creates a new discovery engine using the machine's public identity.
func NewService(id identity.Public) *Service {
	return &Service{
		identity: id,
		peers:    make(map[string]Peer),
	}
}

// Startup begins broadcasting our identity and listening for peers.
func (s *Service) Startup(ctx context.Context) error {
	s.ctx = ctx

	// Encode our identity into TXT records for peers to read
	txtRecords := []string{
		"node_id=" + s.identity.NodeID,
		"uuid=" + s.identity.UUID,
		"public_key=" + s.identity.PublicKey,
	}

	// For now, assume transport engine will bind on port 4242.
	// This will be dynamic once Engine 03 is implemented.
	var err error
	s.server, err = zeroconf.Register("Nodex-"+s.identity.NodeID, "_nodex._tcp", "local.", 4242, txtRecords, nil)
	if err != nil {
		return err
	}

	go s.browse()

	return nil
}

// browse continually scans the local network for other _nodex._tcp services.
func (s *Service) browse() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Println("discovery: failed to initialize resolver:", err)
		return
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			// Skip our own broadcast
			if entry.Instance == "Nodex-"+s.identity.NodeID {
				continue
			}

			peer := Peer{
				Port: entry.Port,
			}

			// Parse identity out of the TXT records
			for _, txt := range entry.Text {
				if len(txt) > 8 && txt[:8] == "node_id=" {
					peer.NodeID = txt[8:]
				} else if len(txt) > 5 && txt[:5] == "uuid=" {
					peer.UUID = txt[5:]
				} else if len(txt) > 11 && txt[:11] == "public_key=" {
					peer.PublicKey = txt[11:]
				}
			}

			var ips []string
			for _, ip := range entry.AddrIPv4 {
				ips = append(ips, ip.String())
			}
			for _, ip := range entry.AddrIPv6 {
				ips = append(ips, ip.String())
			}
			peer.IPs = ips

			// Only track valid peers. We disambiguate strictly by UUID as required.
			if peer.UUID != "" {
				s.peersMu.Lock()
				s.peers[peer.UUID] = peer
				s.peersMu.Unlock()

				if s.ctx != nil {
					// Notify the frontend that a new peer was found or updated
					runtime.EventsEmit(s.ctx, "peer:discovered", peer)
				}
			}
		}
	}(entries)

	// Continually browse, refreshing the query every 15 seconds
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		err = resolver.Browse(ctx, "_nodex._tcp", "local.", entries)
		if err != nil {
			log.Println("discovery: browse error:", err)
		}
		<-ctx.Done()
		cancel()
	}
}

// Shutdown stops the mDNS broadcaster.
func (s *Service) Shutdown(ctx context.Context) {
	if s.server != nil {
		s.server.Shutdown()
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
