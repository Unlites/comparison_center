package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Unlites/comparison_center/backend/config"
	ch "github.com/Unlites/comparison_center/backend/internal/comparison/delivery/http/v1"
	cr "github.com/Unlites/comparison_center/backend/internal/comparison/repository"
	cu "github.com/Unlites/comparison_center/backend/internal/comparison/usecase"
	coh "github.com/Unlites/comparison_center/backend/internal/customoption/delivery/http/v1"
	cor "github.com/Unlites/comparison_center/backend/internal/customoption/repository"
	cou "github.com/Unlites/comparison_center/backend/internal/customoption/usecase"
	middleware "github.com/Unlites/comparison_center/backend/internal/middleware/http"
	oh "github.com/Unlites/comparison_center/backend/internal/object/delivery/http/v1"
	or "github.com/Unlites/comparison_center/backend/internal/object/repository"
	ou "github.com/Unlites/comparison_center/backend/internal/object/usecase"
	ocor "github.com/Unlites/comparison_center/backend/internal/object_customoption/repository"
	g "github.com/Unlites/comparison_center/backend/pkg/generator"
	r "github.com/Unlites/comparison_center/backend/pkg/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	log    *slog.Logger
	cfg    *config.Config
	server *http.Server
}

func NewApp(ctx context.Context, log *slog.Logger, cfg *config.Config) (*App, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
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
	objectUsecase := ou.NewObjectUsecase(objectRepository, objectCustomOptionRepository)
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

	return &App{
		log:    log,
		cfg:    cfg,
		server: srv,
	}, nil
}

func (a *App) Run() error {
	if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
