package handlers

import (
	"seagle/core/services"
)

type GetConfigHandler struct {
	configService *services.ConfigService
}

func NewGetConfigHandler(configService *services.ConfigService) *GetConfigHandler {
	return &GetConfigHandler{
		configService: configService,
	}
}

type GetConfigOutput struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Config  AppConfig `json:"config"`
}

type AppConfig struct {
	OpenAIAPIKey string `json:"openAIAPIKey"`
}

func (h *GetConfigHandler) GetConfig() (*GetConfigOutput, error) {
	cfg, err := h.configService.GetConfig()
	if err != nil {
		return &GetConfigOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &GetConfigOutput{
		Success: true,
		Message: "Configuration updated successfully",
		Config: AppConfig{
			OpenAIAPIKey: cfg.OpenAIAPIKey,
		},
	}, nil
}
