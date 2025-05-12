package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Host   string `env:"HOST_IP_ADDRESS" env-default:"localhost"`
	Port   string `env:"HOST_PORT" env-default:"5005"`
	Level  string `env:"LOG_LEVEL" env-default:"info"`
	DbName string `yaml:"name" env-default:"dreams.sqlite"`
}

// NewConfig returns app config.
func NewConfig(logger *zap.Logger) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		logger.Fatal("config error", zap.Error(err))
		return nil, err
	}

	return cfg, nil
}
