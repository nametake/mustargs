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

	for _, rule := range config.Rules {
		for _, arg := range rule.Args {
			if arg.PkgName == "" {
				arg.PkgName = extractPkgName(arg.Pkg)
			}
		}
	}

	return config, nil
}
