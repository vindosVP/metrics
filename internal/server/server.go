package server

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/handlers"
	"github.com/vindosVP/metrics/internal/middleware"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/server/loader"
	"github.com/vindosVP/metrics/internal/storage/filestorage"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	UpdateGauge(name string, v float64) (float64, error)
	UpdateCounter(name string, v int64) (int64, error)
	SetCounter(name string, v int64) (int64, error)
	GetGauge(name string) (float64, error)
	GetAllGauge() (map[string]float64, error)
	GetCounter(name string) (int64, error)
	GetAllCounter() (map[string]int64, error)
}

func Run(cfg *config.ServerConfig) error {

	gRepo := repos.NewGaugeRepo()
	cRepo := repos.NewCounterRepo()
	var storage MetricsStorage
	if cfg.StoreInterval != time.Duration(0) {
		storage = memstorage.New(gRepo, cRepo)
	} else {
		storage = filestorage.NewFileStorage(gRepo, cRepo, cfg.FileStoragePath)
	}

	if cfg.Restore {
		dumpLoader := loader.New(cfg.FileStoragePath, storage)
		err := dumpLoader.LoadMetrics()
		if err != nil {
			logger.Log.Error("Failed to load dump", zap.Error(err))
		}
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger, middleware.Decompress, chiMiddleware.Compress(5))
	r.Post("/update/", handlers.UpdateBody(storage))
	r.Post("/value/", handlers.GetBody(storage))
	r.Post("/update/{type}/{name}/{value}", handlers.Update(storage))
	r.Get("/value/{type}/{name}", handlers.Get(storage))
	r.Get("/", handlers.List(storage))
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	if cfg.StoreInterval != time.Duration(0) {
		svr := filestorage.NewSaver(cfg.FileStoragePath, cfg.StoreInterval, storage)
		go svr.Run()
		defer svr.Stop()
	}

	log.Printf("Running server on %s", cfg.RunAddr)
	err := http.ListenAndServe(cfg.RunAddr, r)
	if err != nil {
		return err
	}

	return nil
}
