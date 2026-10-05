package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"nodex-desktop/engine/discovery"
	"nodex-desktop/engine/identity"
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

	// Initialize discovery engine
	nodeDiscovery := discovery.NewService(nodeIdentity.GetIdentity())

	// Create an instance of the app structure
	app := NewApp(nodeIdentity, nodeDiscovery)

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
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
