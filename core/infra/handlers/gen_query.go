package handlers

import (
	"seagle/core/services"
	"seagle/core/services/types"
)

// GenerateQueryInput represents the input for the GenerateQuery handler
type GenerateQueryInput struct {
	Database string `json:"database"`
	Prompt   string `json:"prompt"`
}

// GenerateQueryOutput represents the output for the GenerateQuery handler
type GenerateQueryOutput struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message,omitempty"`
	Result  *types.GenerateQueryResult `json:"result,omitempty"`
}

// GenQueryHandler handles AI query generation requests
type GenQueryHandler struct {
	connectionService *services.ConnectionService
}

// NewGenQueryHandler creates a new GenQueryHandler instance
func NewGenQueryHandler(connectionService *services.ConnectionService) *GenQueryHandler {
	return &GenQueryHandler{
		connectionService: connectionService,
	}
}

// GenerateQuery processes the query generation request
func (h *GenQueryHandler) GenerateQuery(input GenerateQueryInput) (*GenerateQueryOutput, error) {
	if input.Prompt == "" {
		return &GenerateQueryOutput{
			Success: false,
			Message: "Prompt cannot be empty",
		}, nil
	}

	if input.Database == "" {
		return &GenerateQueryOutput{
			Success: false,
			Message: "Database must be specified",
		}, nil
	}

	request := types.GenerateQueryRequest{
		Database: input.Database,
		Prompt:   input.Prompt,
	}

	result, err := h.connectionService.GenerateQuery(request)
	if err != nil {
		return &GenerateQueryOutput{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &GenerateQueryOutput{
		Success: true,
		Message: "Query generated successfully",
		Result:  result,
	}, nil
}
