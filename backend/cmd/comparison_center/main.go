package main

import (
	"context"
	"log"

	"github.com/Unlites/backend/comparison_center/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	ctx := context.Background()

	log.Fatal(cfg.MongoURI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

}
