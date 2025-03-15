package config

import (
	"fmt"
	"os"

	"github.com/pycnick/cagenerator/internal/types"
	"gopkg.in/yaml.v2"
)

// Config represents the main configuration structure
type Config struct {
	Entities []types.Entity `yaml:"entities"`
}

// LoadConfig loads and validates the configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Entities) == 0 {
		return fmt.Errorf("at least one entity is required")
	}
	for i, entity := range c.Entities {
		if entity.Name == "" {
			return fmt.Errorf("entity name is required for entity at index %d", i)
		}
	}
	return nil
}
