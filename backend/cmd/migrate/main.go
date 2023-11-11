package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Unlites/comparison_center/backend/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	operation := os.Getenv("MIGRATE_OPERATION")
	if operation == "" {
		slog.Error("MIGRATE_OPERATION is not set. Must be 'up' or 'down'")
		os.Exit(1)
	}

	if operation != "up" && operation != "down" {
		slog.Error("invalid MIGRATE_OPERATION. Must be 'up' or 'down'")
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to init config", "detail", err)
		os.Exit(1)
	}

	m, err := migrate.New(
		fmt.Sprintf("file:///%s", cfg.DB.MigrationsDir),
		cfg.DB.MongoURI,
	)

	if err != nil {
		slog.Error("failed to init migrate", "detail", err)
		os.Exit(1)
	}

	if operation == "up" {
		if err := m.Up(); err != nil {
			slog.Error("failed to migrate up", "detail", err)
			os.Exit(1)
		}
	} else {
		if err := m.Down(); err != nil {
			slog.Error("failed to migrate down", "detail", err)
			os.Exit(1)
		}
	}

	slog.Info("migration completed")
}
