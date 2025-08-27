package main

import (
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"seagle/core/domain"
	"seagle/core/infra/handlers"
	"seagle/core/infra/persistence"
	"seagle/core/services"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create an instance of the app structure
	app := NewApp()

	serviceFactory := domain.NewServiceFactory()
	metadataFactory := domain.NewMetadataFactory(serviceFactory)

	connectionRepo := persistence.NewConnection(persistence.FileAtHomeDir(".seagle", "data", "connections.json"))
	metadataRepo := persistence.NewMetadataRepository(persistence.FileAtHomeDir(".seagle", "data", "metadata.json"))

	// Get OpenAI API key from environment variable
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	openaiClient := services.NewOpenAIClient(openaiAPIKey)

	connectionService := services.NewConnectionService(connectionRepo, metadataRepo, serviceFactory, metadataFactory, openaiClient)

	connectHnd := handlers.NewConnectHandler(connectionService)
	testConnHnd := handlers.NewTestConnectionHandler(connectionService)
	disconnectHnd := handlers.NewDisconnectHandler(connectionService)
	getTablesHnd := handlers.NewGetTablesHandler(connectionService)
	getTableColumnsHnd := handlers.NewGetTableColumnsHandler(connectionService)
	executeQueryHnd := handlers.NewExecuteQueryHandler(connectionService)
	listConnHnd := handlers.NewListConnectionsHandler(connectionService)
	connectByIDHnd := handlers.NewConnectByIDHandler(connectionService)
	analyzeMetadataHnd := handlers.NewAnalyzeMetadataHandler(connectionService)
	genQueryHnd := handlers.NewGenQueryHandler(connectionService)

	// Create application with options
	err = wails.Run(&options.App{
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
			analyzeMetadataHnd,
			genQueryHnd,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
