package handlers

import (
	"seagle/core/services"
)

type ConnectByIDHandler struct {
	connectionService *services.ConnectionService
}

func NewConnectByIDHandler(connectionService *services.ConnectionService) *ConnectByIDHandler {
	return &ConnectByIDHandler{
		connectionService: connectionService,
	}
}

type ConnectByIDInput struct {
	ID string `json:"id"`
}

type ConnectByIDOutput struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message"`
	Databases []string `json:"databases,omitempty"`
}

func (h *ConnectByIDHandler) ConnectByID(input ConnectByIDInput) (*ConnectByIDOutput, error) {
	if input.ID == "" {
		return &ConnectByIDOutput{
			Success: false,
			Message: "Connection ID is required",
		}, nil
	}

	_, err := h.connectionService.ConnectByID(input.ID)
	if err != nil {
		return &ConnectByIDOutput{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Get list of databases after successful connection
	databases, err := h.connectionService.GetDatabases(input.ID)
	if err != nil {
		return &ConnectByIDOutput{
			Success: false,
			Message: "Connected successfully but failed to retrieve databases: " + err.Error(),
		}, nil
	}

	return &ConnectByIDOutput{
		Success:   true,
		Message:   "Connected successfully",
		Databases: databases,
	}, nil
}
