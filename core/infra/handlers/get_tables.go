package handlers

import (
	"seagle/core/services"
)

// GetTablesInput represents the input for the GetTables handler
type GetTablesInput struct {
	ID       string `json:"id"`
	Database string `json:"database"`
}

// GetTablesOutput represents the output for the GetTables handler
type GetTablesOutput struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Tables  []string `json:"tables,omitempty"`
}

// GetTablesHandler handles table listing requests
type GetTablesHandler struct {
	connectionService *services.ConnectionService
}

// NewGetTablesHandler creates a new GetTablesHandler instance
func NewGetTablesHandler(connectionService *services.ConnectionService) *GetTablesHandler {
	return &GetTablesHandler{
		connectionService: connectionService,
	}
}

// GetTables processes the table listing request
func (h *GetTablesHandler) GetTables(input GetTablesInput) (*GetTablesOutput, error) {
	tables, err := h.connectionService.GetTables(input.ID, input.Database)
	if err != nil {
		return &GetTablesOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &GetTablesOutput{
		Success: true,
		Message: "Tables retrieved successfully",
		Tables:  tables,
	}, nil
}
