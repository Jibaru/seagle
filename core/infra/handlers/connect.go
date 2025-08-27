package handlers

import (
	"seagle/core/services"
	"seagle/core/services/types"
)

// ConnectInput represents the input for the Connect handler
type ConnectInput struct {
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Database            string `json:"database"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	SSLMode             string `json:"sslmode"`
	ConnectionString    string `json:"connectionString"`
	UseConnectionString bool   `json:"useConnectionString"`
}

// ConnectOutput represents the output for the Connect handler
type ConnectOutput struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message,omitempty"`
	Databases []string `json:"databases,omitempty"`
	ID        string   `json:"id"`
}

// ConnectHandler handles database connection requests
type ConnectHandler struct {
	connectionService *services.ConnectionService
}

// NewConnectHandler creates a new ConnectHandler instance
func NewConnectHandler(connectionService *services.ConnectionService) *ConnectHandler {
	return &ConnectHandler{
		connectionService: connectionService,
	}
}

// Connect processes the connection request
func (h *ConnectHandler) Connect(input ConnectInput) (*ConnectOutput, error) {
	res, err := h.connectionService.Connect(types.DatabaseConfig{
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
		return &ConnectOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Fetch databases after successful connection
	databases, err := h.connectionService.GetDatabases(res.ID)
	if err != nil {
		// Connection succeeded but database fetching failed - still return success
		return &ConnectOutput{
			Success: true,
			Message: "Connected successfully, but failed to fetch databases: " + err.Error(),
		}, nil
	}

	return &ConnectOutput{
		Success:   true,
		Message:   "Connected successfully",
		Databases: databases,
		ID:        res.ID,
	}, nil
}
