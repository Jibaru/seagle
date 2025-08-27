package handlers

import (
	"seagle/core/services"
)

type SetConfigHandler struct {
	configService *services.ConfigService
}

func NewSetConfigHandler(configService *services.ConfigService) *SetConfigHandler {
	return &SetConfigHandler{
		configService: configService,
	}
}

type SetConfigInput struct {
	OpenAIAPIKey string `json:"openAIAPIKey"`
}

type SetConfigOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *SetConfigHandler) SetConfig(input SetConfigInput) (*SetConfigOutput, error) {
	err := h.configService.SetConfig(input.OpenAIAPIKey)
	if err != nil {
		return &SetConfigOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &SetConfigOutput{
		Success: true,
		Message: "Configuration updated successfully",
	}, nil
}
