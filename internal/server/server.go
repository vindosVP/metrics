// Package server is used to start the http-server to collect metrics
package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/server/grpcserver"
	"github.com/vindosVP/metrics/internal/server/httpserver"
	"github.com/vindosVP/metrics/internal/server/loader"
	"github.com/vindosVP/metrics/internal/storage/dbstorage"
	"github.com/vindosVP/metrics/internal/storage/filestorage"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"github.com/vindosVP/metrics/pkg/logger"
)

// MetricsStorage consists methods to save and get data from the storage
//
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

type pServer interface {
	Run(wg *sync.WaitGroup)
	Stop(wg *sync.WaitGroup)
}

type Server struct {
	http pServer
	grpc pServer
}

func (s *Server) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Add(1)
	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sig
		logger.Log.Info("Got stop signal, stopping")
		go s.http.Stop(wg)
		go s.grpc.Stop(wg)
	}()

	go s.http.Run(wg)
	go s.grpc.Run(wg)

	wg.Wait()
	logger.Log.Info("Server stopped")
}

func withHTTPServer(hs pServer) func(*Server) {
	return func(s *Server) {
		s.http = hs
	}
}

func withGRPCServer(hs pServer) func(*Server) {
	return func(s *Server) {
		s.grpc = hs
	}
}

func newServer(opts ...func(*Server)) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func New(cfg *config.ServerConfig) (*Server, error) {
	s, err := storage(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	hs, err := httpserver.New(s, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	gs, err := grpcserver.New(s, cfg.RPCAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	return newServer(withHTTPServer(hs), withGRPCServer(gs)), nil
}

func storage(cfg *config.ServerConfig) (MetricsStorage, error) {
	if cfg.DatabaseDNS != "" {
		s, err := dbStorage(cfg.DatabaseDNS)
		if err != nil {
			return nil, fmt.Errorf("failed to create database storage: %w", err)
		}
		return s, nil
	} else {
		s, err := memStorage(cfg.StoreInterval, cfg.Restore, cfg.FileStoragePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create inmemory storage: %w", err)
		}
		return s, nil
	}
}

func dbStorage(dsn string) (MetricsStorage, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to databse: %w", err)
	}
	logger.Log.Info("Creating tables")
	err = createTables(pool)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	logger.Log.Info("Tables created successfully")
	return dbstorage.New(pool), nil
}

func memStorage(si time.Duration, restore bool, dump string) (MetricsStorage, error) {

	var s MetricsStorage

	gRepo := repos.NewGaugeRepo()
	cRepo := repos.NewCounterRepo()
	if si != time.Duration(0) {
		logger.Log.Info("Starting saver")
		s = memstorage.New(gRepo, cRepo)
		svr := filestorage.NewSaver(dump, si, s)
		go svr.Run()
	} else {
		s = filestorage.NewFileStorage(gRepo, cRepo, dump)
	}
	if restore {
		dumpLoader := loader.New(dump, s)
		err := dumpLoader.LoadMetrics()
		if err != nil {
			return nil, fmt.Errorf("failed to load dump: %w", err)
		}
	}
	return s, nil
}

func createTables(pool *pgxpool.Pool) error {
	ctx := context.Background()
	query := `CREATE TABLE IF NOT EXISTS gauges (id TEXT NOT NULL PRIMARY KEY, value DOUBLE PRECISION NOT NULL);
			  CREATE TABLE IF NOT EXISTS counters (id TEXT NOT NULL PRIMARY KEY, value BIGINT NOT NULL)`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
