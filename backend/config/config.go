package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type DB struct {
	MongoURI string `env:"MONGODB_URI" env-required:"true"`
}

type Config struct {
	DB
}

func NewConfig() (*Config, error) {
	var cfg Config

	cfgPath := filepath.Join(os.Getenv("CONFIG_PATH"))
	if cfgPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, fmt.Errorf("error while reading config: %w", err)
	}

	return &cfg, nil
}
