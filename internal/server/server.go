// Package server is used to start the http-server to collect metrics
package server

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/handlers"
	"github.com/vindosVP/metrics/internal/middleware"
	"github.com/vindosVP/metrics/internal/models"
	pb "github.com/vindosVP/metrics/internal/proto"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/server/loader"
	"github.com/vindosVP/metrics/internal/service"
	"github.com/vindosVP/metrics/internal/storage/dbstorage"
	"github.com/vindosVP/metrics/internal/storage/filestorage"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"github.com/vindosVP/metrics/pkg/encryption"
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

// Run starts the http server
func Run(cfg *config.ServerConfig) error {
	useDatabase := cfg.DatabaseDNS != ""
	var mux *chi.Mux
	var grpcServer *grpc.Server
	if useDatabase {
		logger.Log.Info("Starting database server")
		logger.Log.Info("Connecting to database")
		ctx := context.Background()
		pool, err := pgxpool.New(ctx, cfg.DatabaseDNS)
		if err != nil {
			logger.Log.Error("Failed to connect to database")
			return err
		}
		logger.Log.Info("Connected successfully")
		defer pool.Close()
		dbmux, a, err := setupDBServer(cfg, pool)
		if err != nil {
			return err
		}
		mux = dbmux
		grpcServer = a
	} else {
		memmux, a, err := setupInmemoryServer(cfg)
		if err != nil {
			return err
		}
		mux = memmux
		grpcServer = a
	}

	return startServer(cfg.RunAddr, grpcServer, mux)
}

func startServer(addr string, a *grpc.Server, mux *chi.Mux) error {

	svr := http.Server{Addr: addr, Handler: mux}

	sd := make(chan struct{})
	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sig
		logger.Log.Info("Got stop signal, stopping")
		err := svr.Shutdown(context.Background())
		if err != nil {
			logger.Log.Error("failed to shutdown http server", zap.Error(err))
		}
		close(sd)
	}()

	logger.Log.Info(fmt.Sprintf("Running server on %s", addr))
	err := svr.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	<-sd
	logger.Log.Info("Stopped successfully")

	return nil
}

func setupDBServer(cfg *config.ServerConfig, pool *pgxpool.Pool) (*chi.Mux, *grpc.Server, error) {

	logger.Log.Info("Creating tables")
	err := createTables(pool)
	if err != nil {
		logger.Log.Error("Failed to create tables")
		return nil, nil, err
	}
	logger.Log.Info("Created successfully")
	storage := dbstorage.New(pool)

	var cryptoKey *rsa.PrivateKey = nil
	if cfg.CryptoKeyFile != "" {
		cryptoKey, err = encryption.PrivateKeyFromFile(cfg.CryptoKeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get crypto key: %w", err)
		}
	}

	var subnet *net.IPNet = nil
	if cfg.TrustedSubnet != "" {
		_, trustedNet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse trusted subnet: %w", err)
		}
		subnet = trustedNet
		logger.Log.Info(fmt.Sprintf("Trusted subnet: %s", subnet.String()))
	}

	r := chi.NewRouter()
	r.Use(middleware.Sign(cfg.Key))
	r.Group(func(r chi.Router) { // group with hash validation
		setupMiddlewares(r, cfg.Key, cryptoKey, true, subnet)
		r.Post("/update/", handlers.UpdateBody(storage))
		r.Post("/updates/", handlers.UpdateBatch(storage))
		r.Post("/value/", handlers.GetBody(storage))
	})
	r.Group(func(r chi.Router) {
		setupMiddlewares(r, "", nil, false, subnet)
		r.Post("/update/{type}/{name}/{value}", handlers.Update(storage))
		r.Get("/value/{type}/{name}", handlers.Get(storage))
		r.Get("/", handlers.List(storage))
		r.Get("/ping", handlers.Ping(pool))
	})
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	a := grpc.NewServer()
	pb.RegisterMetricsServer(a, service.NewMetricsServer(storage))

	return r, a, nil
}

func setupInmemoryServer(cfg *config.ServerConfig) (*chi.Mux, *grpc.Server, error) {
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

	var cryptoKey *rsa.PrivateKey = nil
	if cfg.CryptoKeyFile != "" {
		ck, err := encryption.PrivateKeyFromFile(cfg.CryptoKeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get crypto key: %w", err)
		}
		cryptoKey = ck
	}

	var subnet *net.IPNet = nil
	if cfg.TrustedSubnet != "" {
		_, trustedNet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse trusted subnet: %w", err)
		}
		subnet = trustedNet
		logger.Log.Info(fmt.Sprintf("Trusted subnet: %s", subnet.String()))
	}

	r := chi.NewRouter()
	r.Use(middleware.Sign(cfg.Key))
	r.Group(func(r chi.Router) { // group with hash validation
		setupMiddlewares(r, cfg.Key, cryptoKey, true, subnet)
		r.Post("/update/", handlers.UpdateBody(storage))
		r.Post("/updates/", handlers.UpdateBatch(storage))
		r.Post("/value/", handlers.GetBody(storage))
	})
	r.Group(func(r chi.Router) {
		setupMiddlewares(r, "", nil, false, subnet)
		r.Post("/update/{type}/{name}/{value}", handlers.Update(storage))
		r.Get("/value/{type}/{name}", handlers.Get(storage))
		r.Get("/", handlers.List(storage))
	})
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	a := grpc.NewServer()
	pb.RegisterMetricsServer(a, service.NewMetricsServer(storage))

	return r, a, nil
}

func setupMiddlewares(r chi.Router, key string, cryptoKey *rsa.PrivateKey, validateHash bool, tNet *net.IPNet) {
	r.Use(chiMiddleware.Logger)
	if tNet != nil {
		r.Use(middleware.CheckSubnet(*tNet))
	}
	if validateHash {
		r.Use(middleware.ValidateHMAC(key))
	}
	if cryptoKey != nil {
		r.Use(middleware.Decode(cryptoKey))
	}
	r.Use(middleware.Decompress)
	r.Use(chiMiddleware.Compress(5))
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
