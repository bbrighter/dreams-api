package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		DB   DB   `yaml:"db"`
		API  API  `yaml:"api"`
		Logs Logs `yaml:"logs"`
	}

	DB struct {
		Name string `yaml:"name" env-default:"dreams.sqlite"`
	}

	API struct {
		Host string `yaml:"host" env:"HOST_IP_ADDRESS" env-default:"localhost"`
		Port string `yaml:"port" env-required:"true"`
	}

	Logs struct {
		Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
		Folder string `yaml:"folder" env:"LOG_FOLDER"`
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
