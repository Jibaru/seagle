package handlers

import (
	"seagle/core/services"
	"seagle/core/services/types"
)

// GetTableColumnsInput represents the input for the GetTableColumns handler
type GetTableColumnsInput struct {
	ID       string `json:"id"`
	Database string `json:"database"`
	Table    string `json:"table"`
}

// GetTableColumnsOutput represents the output for the GetTableColumns handler
type GetTableColumnsOutput struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Columns []types.TableColumn `json:"columns,omitempty"`
}

// GetTableColumnsHandler handles table column listing requests
type GetTableColumnsHandler struct {
	connectionService *services.ConnectionService
}

// NewGetTableColumnsHandler creates a new GetTableColumnsHandler instance
func NewGetTableColumnsHandler(connectionService *services.ConnectionService) *GetTableColumnsHandler {
	return &GetTableColumnsHandler{
		connectionService: connectionService,
	}
}

// GetTableColumns processes the table column listing request
func (h *GetTableColumnsHandler) GetTableColumns(input GetTableColumnsInput) (*GetTableColumnsOutput, error) {
	columns, err := h.connectionService.GetTableColumns(input.ID, input.Database, input.Table)
	if err != nil {
		return &GetTableColumnsOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &GetTableColumnsOutput{
		Success: true,
		Message: "Columns retrieved successfully",
		Columns: columns,
	}, nil
}
