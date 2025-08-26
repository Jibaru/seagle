package handlers

import (
	"context"

	"seagle/core/services"
	"seagle/core/services/types"
)

type ListConnectionsHandler struct {
	connectionService *services.ConnectionService
}

func NewListConnectionsHandler(connectionService *services.ConnectionService) *ListConnectionsHandler {
	return &ListConnectionsHandler{
		connectionService: connectionService,
	}
}

type ListConnectionsOutput struct {
	Success     bool                        `json:"success"`
	Message     string                      `json:"message"`
	Connections []types.ConnectionSummary   `json:"connections"`
}

func (h *ListConnectionsHandler) ListConnections(ctx context.Context) (*ListConnectionsOutput, error) {
	connections, err := h.connectionService.ListConnections()
	if err != nil {
		return &ListConnectionsOutput{
			Success:     false,
			Message:     err.Error(),
			Connections: []types.ConnectionSummary{},
		}, nil
	}

	return &ListConnectionsOutput{
		Success:     true,
		Message:     "Connections listed successfully",
		Connections: connections,
	}, nil
}
