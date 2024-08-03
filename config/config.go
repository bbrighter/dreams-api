package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		DB  DB  `yaml:"db"`
		API API `yaml:"api"`
	}

	DB struct {
		Name string `yaml:"name" env-required:"true"`
	}

	API struct {
		Host string `yaml:"host" env:"HOST_IP_ADDRESS" env-default:"localhost"`
		Port string `yaml:"port" env-required:"true"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
