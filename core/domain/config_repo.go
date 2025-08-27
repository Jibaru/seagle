package domain

type ConfigRepo interface {
	Save(cfg *Config) error
	Find() (*Config, error)
}
