package domain

type Config struct {
	openAIAPIKey string
}

func NewConfig(openAIAPIKey string) (*Config, error) {
	return &Config{
		openAIAPIKey: openAIAPIKey,
	}, nil
}

func NewConfigFromMap(data map[string]interface{}) *Config {
	openAIAPIKey, _ := data["openAIAPIKey"].(string)
	return &Config{
		openAIAPIKey: openAIAPIKey,
	}
}

func (c *Config) OpenAIAPIKey() string {
	return c.openAIAPIKey
}

func (c *Config) SetOpenAIAPIKey(key string) {
	c.openAIAPIKey = key
}

func (c *Config) ToMap() map[string]any {
	return map[string]any{
		"openAIAPIKey": c.openAIAPIKey,
	}
}
