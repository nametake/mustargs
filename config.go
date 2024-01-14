package mustargs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

func loadConfig(filepath string) (*Config, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var config *Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return config, nil
}
