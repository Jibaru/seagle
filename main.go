package main

import (
	"embed"
	"seagle/core/infra/handlers"
	"seagle/core/infra/persistence"
	"seagle/core/services"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	connectionRepo := persistence.NewConnection("connections.json")

	connectionService := services.NewConnectionService(connectionRepo)

	connectHnd := handlers.NewConnectHandler(connectionService)
	testConnHnd := handlers.NewTestConnectionHandler(connectionService)
	disconnectHnd := handlers.NewDisconnectHandler(connectionService)
	getTablesHnd := handlers.NewGetTablesHandler(connectionService)
	getTableColumnsHnd := handlers.NewGetTableColumnsHandler(connectionService)
	executeQueryHnd := handlers.NewExecuteQueryHandler(connectionService)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "seagle",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
			connectHnd,
			testConnHnd,
			disconnectHnd,
			getTablesHnd,
			getTableColumnsHnd,
			executeQueryHnd,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
