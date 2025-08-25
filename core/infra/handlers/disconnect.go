package handlers

import "seagle/core/services"

// DisconnectHandler handles database disconnection requests
type DisconnectHandler struct {
	connectionService *services.ConnectionService
}

// NewDisconnectHandler creates a new DisconnectHandler instance
func NewDisconnectHandler(connectionService *services.ConnectionService) *DisconnectHandler {
	return &DisconnectHandler{
		connectionService: connectionService,
	}
}

// Disconnect processes the disconnect request
func (h *DisconnectHandler) Disconnect() error {
	err := h.connectionService.Disconnect()
	if err != nil {
		return err
	}

	return nil
}
