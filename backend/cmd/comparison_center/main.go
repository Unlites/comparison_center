package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Unlites/comparison_center/backend/config"
	ch "github.com/Unlites/comparison_center/backend/internal/adapters/handlers/http/v1/comparison"
	coh "github.com/Unlites/comparison_center/backend/internal/adapters/handlers/http/v1/customoption"
	"github.com/Unlites/comparison_center/backend/internal/adapters/handlers/http/v1/middleware"
	oh "github.com/Unlites/comparison_center/backend/internal/adapters/handlers/http/v1/object"
	cr "github.com/Unlites/comparison_center/backend/internal/adapters/repositories/comparison"
	cor "github.com/Unlites/comparison_center/backend/internal/adapters/repositories/customoption"
	or "github.com/Unlites/comparison_center/backend/internal/adapters/repositories/object"
	ocor "github.com/Unlites/comparison_center/backend/internal/adapters/repositories/object_customoption"
	cu "github.com/Unlites/comparison_center/backend/internal/application/comparison"
	cou "github.com/Unlites/comparison_center/backend/internal/application/customoption"
	ou "github.com/Unlites/comparison_center/backend/internal/application/object"
	g "github.com/Unlites/comparison_center/backend/pkg/generator"
	"github.com/Unlites/comparison_center/backend/pkg/metrics"
	"github.com/Unlites/comparison_center/backend/pkg/parser"
	r "github.com/Unlites/comparison_center/backend/pkg/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to init config", "detail", err)
		os.Exit(1)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parser.ParseSlogLevel(cfg.LogLevel),
	}))

	metrics := metrics.NewMetrics(cfg.MetricsAddress)

	go func() {
		log.Info("starting metrics server", "addr", cfg.MetricsAddress)
		if err := metrics.Run(); err != nil {
			log.Error("failed to start metrics server", "detail", err)
			os.Exit(1)
		}
	}()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.URI))
	if err != nil {
		log.Error("failed to connect to mongodb", "detail", err)
		os.Exit(1)
	}

	generator := g.NewGenerator()

	comparisonRepository := cr.NewComparisonRepositoryMongo(client)
	comparisonUsecase := cu.NewComparisonUsecase(comparisonRepository, generator)
	comparisonHandler := ch.NewComparisonHandler(comparisonUsecase)

	customOptionRepository := cor.NewCustomOptionRepositoryMongo(client)
	customOptionUsecase := cou.NewCustomOptionUsecase(customOptionRepository, generator)
	customOptionHandler := coh.NewCustomOptionHandler(customOptionUsecase)

	objectRepository := or.NewObjectRepositoryMongo(client)
	objectCustomOptionRepository := ocor.NewObjectCustomOptionRepositoryMongo(client)
	objectUsecase := ou.NewObjectUsecase(objectRepository, objectCustomOptionRepository, generator)
	objectHandler := oh.NewObjectHandler(objectUsecase, cfg.PhotosDir, cfg.MaxUploadSizeMB)

	router := r.NewDefaultRouter()
	router.Handler.Use(middleware.Metrics)
	router.RegisterHandlers("v1", map[string]http.Handler{
		"comparisons":    comparisonHandler,
		"custom_options": customOptionHandler,
		"objects":        objectHandler,
	})

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router.Handler,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	go func() {
		log.Info("starting application server", "addr", cfg.HttpServer.Address)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start application server", "detail", err)
			os.Exit(1)
		}
	}()

	notifyCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-notifyCtx.Done()
	log.Info("service gracefully stopping...")

	shutDownCtx, cancel := context.WithTimeout(ctx, cfg.HttpServer.ShutdownTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Shutdown(shutDownCtx); err != nil {
			log.Error("failed to stop application server", "detail", err)
			return
		}
		log.Info("application server stopped")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := metrics.Stop(shutDownCtx); err != nil {
			log.Error("failed to stop metrics server", "detail", err)
			return
		}
		log.Info("metrics server stopped")
	}()

	wg.Wait()
	log.Info("service stopped")
}
