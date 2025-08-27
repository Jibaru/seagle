package handlers

import (
	"seagle/core/services"
)

// AnalyzeMetadataInput represents the input for the AnalyzeMetadata handler
type AnalyzeMetadataInput struct {
	ID string `json:"id"`
}

// AnalyzeMetadataOutput represents the output for the AnalyzeMetadata handler
type AnalyzeMetadataOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// AnalyzeMetadataHandler handles connection metadata analysis requests
type AnalyzeMetadataHandler struct {
	connectionService *services.ConnectionService
}

// NewAnalyzeMetadataHandler creates a new AnalyzeMetadataHandler instance
func NewAnalyzeMetadataHandler(connectionService *services.ConnectionService) *AnalyzeMetadataHandler {
	return &AnalyzeMetadataHandler{
		connectionService: connectionService,
	}
}

// AnalyzeMetadata processes the metadata analysis request
func (h *AnalyzeMetadataHandler) AnalyzeMetadata(input AnalyzeMetadataInput) (*AnalyzeMetadataOutput, error) {
	err := h.connectionService.AnalyzeConnectionMetadata(input.ID)
	if err != nil {
		return &AnalyzeMetadataOutput{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &AnalyzeMetadataOutput{
		Success: true,
		Message: "Metadata analysis and persistence completed successfully",
	}, nil
}
