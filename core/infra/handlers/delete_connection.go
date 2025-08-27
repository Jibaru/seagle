package handlers

import "seagle/core/services"

// DeleteConnectionHandler handles database DeleteConnectionion requests
type DeleteConnectionHandler struct {
	connectionService *services.ConnectionService
}

// DeleteConnectionInput represents the input for the DeleteConnection handler
type DeleteConnectionInput struct {
	ID string `json:"id"`
}

// NewDeleteConnectionHandler creates a new DeleteConnectionHandler instance
func NewDeleteConnectionHandler(connectionService *services.ConnectionService) *DeleteConnectionHandler {
	return &DeleteConnectionHandler{
		connectionService: connectionService,
	}
}

// DeleteConnection processes the DeleteConnection request
func (h *DeleteConnectionHandler) DeleteConnection(input DeleteConnectionInput) error {
	err := h.connectionService.DeleteConnection(input.ID)
	if err != nil {
		return err
	}

	return nil
}
