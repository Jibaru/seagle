package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"seagle/core/infra/handlers"
	"seagle/core/infra/persistence"
	"seagle/core/services"
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
	listConnHnd := handlers.NewListConnectionsHandler(connectionService)
	connectByIDHnd := handlers.NewConnectByIDHandler(connectionService)

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "seagle",
		Width:            1024,
		Height:           768,
		WindowStartState: options.Maximised,
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
			listConnHnd,
			connectByIDHnd,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
