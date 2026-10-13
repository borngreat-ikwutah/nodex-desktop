package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"nodex-desktop/engine/discovery"
	"nodex-desktop/engine/identity"
	"nodex-desktop/engine/media"
	"nodex-desktop/engine/transport"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Load (or generate on first boot) this installation's permanent identity
	nodeIdentity, err := identity.NewService()
	if err != nil {
		println("identity error:", err.Error())
		return
	}

	// Initialize transport engine (TCP Server & Client)
	nodeTransport := transport.NewService(nodeIdentity.GetIdentity(), nodeIdentity.SigningKey())

	// Initialize discovery engine (advertising the transport port)
	nodeDiscovery := discovery.NewService(
		nodeIdentity.GetIdentity(),
		discovery.WithTransportPort(nodeTransport.Port()),
	)
	nodeMedia := media.NewService(nodeIdentity.GetIdentity())

	// Create an instance of the app structure
	app := NewApp(nodeIdentity, nodeDiscovery, nodeTransport, nodeMedia)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "nodex-desktop",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
