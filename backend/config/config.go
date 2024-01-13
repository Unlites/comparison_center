package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address         string        `yaml:"address"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	MaxUploadSizeMB int64         `yaml:"max_upload_size_mb"`
}

type DB struct {
	URI           string `yaml:"uri"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type Config struct {
	HttpServer     `yaml:"http_server"`
	MetricsAddress string `yaml:"metrics_address"`
	DB             `yaml:"db"`
	PhotosDir      string `yaml:"photos_dir"`
	LogLevel       string `yaml:"log_level"`
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
