package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/handlers"
	"github.com/vindosVP/metrics/internal/middleware"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/server/loader"
	"github.com/vindosVP/metrics/internal/storage/dbstorage"
	"github.com/vindosVP/metrics/internal/storage/filestorage"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	UpdateGauge(ctx context.Context, name string, v float64) (float64, error)
	UpdateCounter(ctx context.Context, name string, v int64) (int64, error)
	SetCounter(ctx context.Context, name string, v int64) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
	GetAllGauge(ctx context.Context) (map[string]float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAllCounter(ctx context.Context) (map[string]int64, error)
	InsertBatch(ctx context.Context, batch []*models.Metrics) error
}

func Run(cfg *config.ServerConfig) error {
	useDatabase := cfg.DatabaseDNS != ""
	var mux *chi.Mux
	if useDatabase {
		logger.Log.Info("Starting database server")
		logger.Log.Info("Connecting to database")
		db, err := sql.Open("postgres", cfg.DatabaseDNS)
		if err != nil {
			logger.Log.Error("Failed to connect to database")
			return err
		}
		logger.Log.Info("Connected successfully")
		defer db.Close()
		dbmux, err := setupDBServer(db)
		if err != nil {
			return err
		}
		mux = dbmux
	} else {
		memmux, err := setupInmemoryServer(cfg)
		if err != nil {
			return err
		}
		mux = memmux
	}

	logger.Log.Info(fmt.Sprintf("Running server on %s", cfg.RunAddr))
	err := http.ListenAndServe(cfg.RunAddr, mux)
	if err != nil {
		return err
	}

	return nil
}

func setupDBServer(db *sql.DB) (*chi.Mux, error) {

	logger.Log.Info("Creating tables")
	err := createTables(db)
	if err != nil {
		logger.Log.Error("Failed to create tables")
		return nil, err
	}
	logger.Log.Info("Created successfully")
	storage := dbstorage.New(db)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger, middleware.Decompress, chiMiddleware.Compress(5))
	r.Get("/ping", handlers.Ping(db))
	r.Post("/update/", handlers.UpdateBody(storage))
	r.Post("/updates/", handlers.UpdateBatch(storage))
	r.Post("/value/", handlers.GetBody(storage))
	r.Post("/update/{type}/{name}/{value}", handlers.Update(storage))
	r.Get("/value/{type}/{name}", handlers.Get(storage))
	r.Get("/", handlers.List(storage))
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	return r, nil
}

func setupInmemoryServer(cfg *config.ServerConfig) (*chi.Mux, error) {
	logger.Log.Info("Starting inmemory server")

	var storage MetricsStorage
	gRepo := repos.NewGaugeRepo()
	cRepo := repos.NewCounterRepo()
	if cfg.StoreInterval != time.Duration(0) {
		logger.Log.Info("Starting saver")
		storage = memstorage.New(gRepo, cRepo)
		svr := filestorage.NewSaver(cfg.FileStoragePath, cfg.StoreInterval, storage)
		go svr.Run()
		defer svr.Stop()
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
	r.Post("/updates/", handlers.UpdateBatch(storage))
	r.Post("/value/", handlers.GetBody(storage))
	r.Post("/update/{type}/{name}/{value}", handlers.Update(storage))
	r.Get("/value/{type}/{name}", handlers.Get(storage))
	r.Get("/", handlers.List(storage))
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	return r, nil
}

func createTables(db *sql.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	gaugeRequest := "CREATE TABLE IF NOT EXISTS gauges (id VARCHAR(250) NOT NULL PRIMARY KEY, value double precision NOT NULL)"
	counterRequest := "CREATE TABLE IF NOT EXISTS counters (id VARCHAR(250) NOT NULL PRIMARY KEY, value integer NOT NULL)"

	_, err = tx.ExecContext(ctx, gaugeRequest)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, counterRequest)
	if err != nil {
		return err
	}

	return tx.Commit()
}
