package main

import (
	"context"

	"nodex-desktop/engine/discovery"
	"nodex-desktop/engine/identity"
)

// App is the Wails-bound application surface.
type App struct {
	ctx context.Context
	id  *identity.Service
	dis *discovery.Service
}

// NewApp creates a new App application struct backed by the machine identity.
func NewApp(id *identity.Service, dis *discovery.Service) *App {
	return &App{id: id, dis: dis}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.dis != nil {
		a.dis.Startup(ctx)
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
