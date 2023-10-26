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
	MongoURI string `env:"MONGODB_URI" env-required:"true"`
}

type Config struct {
	HttpServer `yaml:"http_server"`
	DB
	PhotosDir string `env:"PHOTOS_DIR" env-required:"true"`
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
