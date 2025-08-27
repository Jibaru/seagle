package services

import "seagle/core/domain"

type ConfigService struct {
	repo domain.ConfigRepo
}

func NewConfigService(repo domain.ConfigRepo) *ConfigService {
	return &ConfigService{
		repo: repo,
	}
}

func (s *ConfigService) SetConfig(
	openAIAPIKey string,
) error {
	cfg, err := s.repo.Find()
	if err != nil {
		return err
	}

	if cfg == nil {
		cfg, err = domain.NewConfig(openAIAPIKey)
		if err != nil {
			return err
		}
	} else {
		cfg.SetOpenAIAPIKey(openAIAPIKey)
	}

	return s.repo.Save(cfg)
}
