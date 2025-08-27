package persistence

import "seagle/core/domain"

type ConfigRepo struct {
	filename string
}

func NewConfigRepo(filename string) *ConfigRepo {
	return &ConfigRepo{
		filename: filename,
	}
}

func (r *ConfigRepo) Save(cfg *domain.Config) error {
	data := cfg.ToMap()
	return saveDataToFile(r.filename, data)
}

func (r *ConfigRepo) Find() (*domain.Config, error) {
	data, err := loadDataFromFile(r.filename)
	if err != nil {
		return nil, err
	}
	return domain.NewConfigFromMap(data), nil
}
