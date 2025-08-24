package handlers

import (
	"seagle/core/services"
	"seagle/core/services/types"
)

// TestConnectionInput represents the input for the TestConnection handler
type TestConnectionInput struct {
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Database            string `json:"database"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	SSLMode             string `json:"sslmode"`
	ConnectionString    string `json:"connectionString"`
	UseConnectionString bool   `json:"useConnectionString"`
}

// TestConnectionHandler handles database connection testing requests
type TestConnectionHandler struct {
	connectionService *services.ConnectionService
}

// NewTestConnectionHandler creates a new TestConnectionHandler instance
func NewTestConnectionHandler(connectionService *services.ConnectionService) *TestConnectionHandler {
	return &TestConnectionHandler{
		connectionService: connectionService,
	}
}

// TestConnection processes the test connection request
func (h *TestConnectionHandler) TestConnection(input TestConnectionInput) error {
	err := h.connectionService.TestConnection(types.DatabaseConfig{
		Host:                input.Host,
		Port:                input.Port,
		Database:            input.Database,
		Username:            input.Username,
		Password:            input.Password,
		SSLMode:             input.SSLMode,
		ConnectionString:    input.ConnectionString,
		UseConnectionString: input.UseConnectionString,
	})
	if err != nil {
		return err
	}

	return nil
}
