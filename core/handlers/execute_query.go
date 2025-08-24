package handlers

import (
	"seagle/core/services"
	"seagle/core/services/types"
)

// ExecuteQueryInput represents the input for the ExecuteQuery handler
type ExecuteQueryInput struct {
	Database string `json:"database"`
	Query    string `json:"query"`
}

// ExecuteQueryOutput represents the output for the ExecuteQuery handler
type ExecuteQueryOutput struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Result  *types.QueryResult  `json:"result,omitempty"`
}

// ExecuteQueryHandler handles query execution requests
type ExecuteQueryHandler struct {
	connectionService *services.ConnectionService
}

// NewExecuteQueryHandler creates a new ExecuteQueryHandler instance
func NewExecuteQueryHandler(connectionService *services.ConnectionService) *ExecuteQueryHandler {
	return &ExecuteQueryHandler{
		connectionService: connectionService,
	}
}

// ExecuteQuery processes the query execution request
func (h *ExecuteQueryHandler) ExecuteQuery(input ExecuteQueryInput) (*ExecuteQueryOutput, error) {
	if input.Query == "" {
		return &ExecuteQueryOutput{
			Success: false,
			Message: "Query cannot be empty",
		}, nil
	}

	result, err := h.connectionService.ExecuteQuery(input.Database, input.Query)
	if err != nil {
		return &ExecuteQueryOutput{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &ExecuteQueryOutput{
		Success: true,
		Message: "Query executed successfully",
		Result:  result,
	}, nil
}