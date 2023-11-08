package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Unlites/comparison_center/backend/config"
	ch "github.com/Unlites/comparison_center/backend/internal/comparison/delivery/http/v1"
	cr "github.com/Unlites/comparison_center/backend/internal/comparison/repository"
	cu "github.com/Unlites/comparison_center/backend/internal/comparison/usecase"
	coh "github.com/Unlites/comparison_center/backend/internal/customoption/delivery/http/v1"
	cor "github.com/Unlites/comparison_center/backend/internal/customoption/repository"
	cou "github.com/Unlites/comparison_center/backend/internal/customoption/usecase"
	oh "github.com/Unlites/comparison_center/backend/internal/object/delivery/http/v1"
	or "github.com/Unlites/comparison_center/backend/internal/object/repository"
	ou "github.com/Unlites/comparison_center/backend/internal/object/usecase"
	ocor "github.com/Unlites/comparison_center/backend/internal/object_customoption/repository"
	r "github.com/Unlites/comparison_center/backend/pkg/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.ErrorContext(ctx, "failed to init config", "detail", err)
		os.Exit(1)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.MongoURI))
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to mongodb", "detail", err)
		os.Exit(1)
	}

	if err := client.Ping(ctx, nil); err != nil {
		slog.ErrorContext(ctx, "failed to ping mongo", "detail", err)
		os.Exit(1)
	}

	comparisonRepository := cr.NewComparisonRepositoryMongo(client)
	comparisonUsecase := cu.NewComparisonUsecase(comparisonRepository)
	comparisonHandler := ch.NewComparisonHandler(comparisonUsecase)

	customOptionRepository := cor.NewCustomOptionRepositoryMongo(client)
	customOptionUsecase := cou.NewCustomOptionUsecase(customOptionRepository)
	customOptionHandler := coh.NewCustomOptionHandler(customOptionUsecase)

	objectRepository := or.NewObjectRepositoryMongo(client)
	objectCustomOptionRepository := ocor.NewObjectCustomOptionRepositoryMongo(client)
	objectUsecase := ou.NewObjectUsecase(objectRepository, objectCustomOptionRepository)
	objectHandler := oh.NewObjectHandler(objectUsecase, cfg.PhotosDir, cfg.MaxUploadSizeMB)

	router := r.NewDefaultRouter()
	router.RegisterHandlers("v1", map[string]http.Handler{
		"comparisons":    comparisonHandler.Router,
		"custom_options": customOptionHandler.Router,
		"objects":        objectHandler.Router,
	})

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router.Handler,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	slog.InfoContext(ctx, "starting server...", "addr", cfg.HttpServer.Address)

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "failed to start server", "detail", err)
			os.Exit(1)
		}
	}()

	slog.InfoContext(ctx, "server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	slog.InfoContext(ctx, "server gracefully stopping...")

	shutDownCtx, cancel := context.WithTimeout(ctx, cfg.HttpServer.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		slog.ErrorContext(ctx, "failed to gracefully stop the server")
		os.Exit(1)
	}

	slog.InfoContext(ctx, "server gracefully stopped")

}
