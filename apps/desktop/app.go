package main

import (
	"context"
	"errors"
	"net"
	"strings"

	"nodex-desktop/engine/discovery"
	"nodex-desktop/engine/identity"
	"nodex-desktop/engine/media"
	"nodex-desktop/engine/transport"
)

// App is the Wails-bound application surface.
type App struct {
	ctx   context.Context
	id    *identity.Service
	dis   *discovery.Service
	tran  *transport.Service
	media *media.Service
}

// NewApp creates a new App application struct backed by the machine identity and engines.
func NewApp(id *identity.Service, dis *discovery.Service, tran *transport.Service, mediaService *media.Service) *App {
	return &App{id: id, dis: dis, tran: tran, media: mediaService}
}

// CreateWebRTCOffer starts a media negotiation and returns the SDP signal.
func (a *App) CreateWebRTCOffer(sessionID string) (media.Signal, error) {
	if a.media == nil {
		return media.Signal{}, errors.New("media service not initialized")
	}
	return a.media.CreateOffer(context.Background(), sessionID)
}

// AcceptWebRTCAnswer completes a media negotiation with a remote SDP answer.
func (a *App) AcceptWebRTCAnswer(signal media.Signal) error {
	if a.media == nil {
		return errors.New("media service not initialized")
	}
	return a.media.AcceptAnswer(context.Background(), signal)
}

// AddWebRTCICE adds a remote ICE candidate to a media session.
func (a *App) AddWebRTCICE(signal media.Signal) error {
	if a.media == nil {
		return errors.New("media service not initialized")
	}
	return a.media.AddICECandidate(context.Background(), signal)
}

// CloseWebRTCSession releases a media session.
func (a *App) CloseWebRTCSession(sessionID string) error {
	if a.media == nil {
		return nil
	}
	return a.media.Close(sessionID)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.dis != nil {
		_ = a.dis.Startup(ctx)
	}
	if a.tran != nil {
		_ = a.tran.Startup(ctx)
	}
}

// shutdown is called when the app terminates.
func (a *App) shutdown(ctx context.Context) {
	if a.dis != nil {
		a.dis.Shutdown(ctx)
	}
	if a.tran != nil {
		_ = a.tran.Shutdown(ctx)
	}
}

// GetNodeID returns this installation's permanent, human readable Nodex ID.
func (a *App) GetNodeID() string {
	if a.id == nil {
		return ""
	}
	return a.id.GetNodeID()
}

// GetIdentity returns the full public identity payload for this installation.
func (a *App) GetIdentity() identity.Public {
	if a.id == nil {
		return identity.Public{}
	}
	return a.id.GetIdentity()
}

// GetPeers returns all currently discovered peers on the LAN.
func (a *App) GetPeers() []discovery.Peer {
	if a.dis == nil {
		return nil
	}
	return a.dis.GetPeers()
}

// FindNode locates a peer on the local network given its 9-digit Node ID.
func (a *App) FindNode(targetID string) (*discovery.Peer, error) {
	if a.dis == nil {
		return nil, discovery.ErrNodeNotFound
	}
	return a.dis.FindNode(a.ctx, targetID)
}

// PingNode sends a TCP ping to the remote address (e.g. "192.168.1.50:4242") and returns latency metrics.
func (a *App) PingNode(addr string) (*transport.PingResult, error) {
	if a.tran == nil {
		return nil, errors.New("transport service not initialized")
	}
	// Pin the TLS certificate when the address belongs to a peer learned by discovery.
	if a.dis != nil {
		if host, _, err := net.SplitHostPort(addr); err == nil {
			for _, peer := range a.dis.GetPeers() {
				for _, ip := range peer.IPs {
					if strings.TrimSpace(ip) == host {
						return a.tran.PingVerified(a.ctx, addr, identity.Public{
							NodeID: peer.NodeID, UUID: peer.UUID, PublicKey: peer.PublicKey,
						})
					}
				}
			}
		}
	}
	return a.tran.Ping(addr)
}

// GetTransportPort returns the local TCP server port for peer connections.
func (a *App) GetTransportPort() int {
	if a.tran == nil {
		return 0
	}
	return a.tran.Port()
}
