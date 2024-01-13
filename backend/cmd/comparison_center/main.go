package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Unlites/comparison_center/backend/config"
	"github.com/Unlites/comparison_center/backend/internal/app"
	"github.com/Unlites/comparison_center/backend/internal/utils/parser"

	"github.com/Unlites/comparison_center/backend/pkg/metrics"
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

	application, err := app.NewApp(ctx, log, cfg)
	if err != nil {
		log.Error("failed to init application", "detail", err)
		os.Exit(1)
	}

	go func() {
		log.Info("starting application server", "addr", cfg.HttpServer.Address)
		if err := application.Run(); err != nil {
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
		if err := application.Stop(shutDownCtx); err != nil {
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

	gracefullyDoneCh := make(chan struct{})

	go func() {
		select {
		case <-shutDownCtx.Done():
			log.Info("service stopped due to shutdown timeout")
		case <-gracefullyDoneCh:
			log.Info("service stopped gracefully")
		}
	}()

	wg.Wait()
	close(gracefullyDoneCh)
}
