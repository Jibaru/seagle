package handlers

import "seagle/core/services"

// DisconnectHandler handles database disconnection requests
type DisconnectHandler struct {
	connectionService *services.ConnectionService
}

// DisconnectInput represents the input for the Disconnect handler
type DisconnectInput struct {
	ID string `json:"id"`
}

// NewDisconnectHandler creates a new DisconnectHandler instance
func NewDisconnectHandler(connectionService *services.ConnectionService) *DisconnectHandler {
	return &DisconnectHandler{
		connectionService: connectionService,
	}
}

// Disconnect processes the disconnect request
func (h *DisconnectHandler) Disconnect(input DisconnectInput) error {
	err := h.connectionService.Disconnect(input.ID)
	if err != nil {
		return err
	}

	return nil
}
