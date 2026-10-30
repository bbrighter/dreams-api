package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Host          string `env:"HOST_IP_ADDRESS" env-default:"localhost"`
	Port          string `env:"HOST_PORT" env-default:"5000"`
	Level         string `env:"LOG_LEVEL" env-default:"info"`
	DbName        string `yaml:"name" env-default:"dreams.sqlite"`
	DataFolder    string `env:"DATA_FOLDER" env-default:"data"`
	Transcription TranscriptionConfig
}

type TranscriptionConfig struct {
	Url      string `env:"TRANSCRIPTION_URL" env-default:"http://localhost:8000"`
	Model    string `env:"TRANSCRIPTION_MODEL" env-default:"Systran/faster-whisper-tiny"`
	Language string `env:"TRANSCRIPTION_LANGUAGE" env-default:"de"`
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
