package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Unlites/backend/comparison_center/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to init config", "detail", err)
		os.Exit(1)
	}

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		slog.Error("failed to connect to mongodb", "detail", err)
		os.Exit(1)
	}

	if err := client.Ping(ctx, nil); err != nil {
		slog.Error(err.Error())
	}

}
