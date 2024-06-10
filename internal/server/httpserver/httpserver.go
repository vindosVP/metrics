package httpserver

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	chiMws "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/handlers"
	"github.com/vindosVP/metrics/internal/middleware"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/encryption"
	"github.com/vindosVP/metrics/pkg/logger"
)

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

type HTTPServer struct {
	s *http.Server
}

func (h *HTTPServer) Run(wg *sync.WaitGroup) {
	err := h.s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		wg.Done()
		logger.Log.Error("failed to start HTTP server", zap.Error(err))
	}
}

func (h *HTTPServer) Stop(wg *sync.WaitGroup) {
	err := h.s.Shutdown(context.Background())
	if err != nil {
		logger.Log.Error("failed to stop HTTP server", zap.Error(err))
	}
	wg.Done()
}

func newHTTPServer(s *http.Server) *HTTPServer {
	return &HTTPServer{s: s}
}

type server struct {
	mux  *chi.Mux
	addr string
}

func newServer(opts ...func(*server)) *http.Server {
	s := &server{
		mux:  chi.NewRouter(),
		addr: "",
	}
	for _, opt := range opts {
		opt(s)
	}
	s.mux.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	return &http.Server{Addr: s.addr, Handler: s.mux}
}

func configuration(c *httpServerConfig) []func(*server) {
	return []func(*server){
		withAddr(c.Addr),
		withMw(chiMws.Logger),
		withMw(middleware.Sign(c.Key)),
		withRouteGroup(legacyGroup(c.Storage, c.Subnet)),
		withRouteGroup(group(c.Storage, c.Key, c.PKey, c.Subnet)),
	}
}

func withAddr(addr string) func(*server) {
	return func(s *server) {
		logger.Log.Info(fmt.Sprintf("Http server listening on %s", addr))
		s.addr = addr
	}
}

func withMw(mw func(next http.Handler) http.Handler) func(*server) {
	return func(s *server) {
		s.mux.Use(mw)
	}
}

func withRouteGroup(fun func(r chi.Router)) func(*server) {
	return func(s *server) {
		s.mux.Group(fun)
	}
}

func legacyGroup(st MetricsStorage, subnet *net.IPNet) func(r chi.Router) {
	return func(r chi.Router) {
		if subnet != nil {
			r.Use(middleware.CheckSubnet(*subnet))
		}
		r.Use(middleware.Decompress)
		r.Use(chiMws.Compress(5))
		r.Post("/update/{type}/{name}/{value}", handlers.Update(st))
		r.Get("/value/{type}/{name}", handlers.Get(st))
		r.Get("/", handlers.List(st))
	}
}

func group(st MetricsStorage, key string, pKey *rsa.PrivateKey, subnet *net.IPNet) func(r chi.Router) {
	return func(r chi.Router) {
		if subnet != nil {
			r.Use(middleware.CheckSubnet(*subnet))
		}
		r.Use(middleware.ValidateHMAC(key))
		if pKey != nil {
			r.Use(middleware.Decode(pKey))
		}
		r.Use(middleware.Decompress)
		r.Use(chiMws.Compress(5))
		r.Post("/update/", handlers.UpdateBody(st))
		r.Post("/updates/", handlers.UpdateBatch(st))
		r.Post("/value/", handlers.GetBody(st))
	}
}

type httpServerConfig struct {
	Subnet  *net.IPNet
	Key     string
	PKey    *rsa.PrivateKey
	Addr    string
	Storage MetricsStorage
}

func newConfig(st MetricsStorage, cfg *config.ServerConfig) (*httpServerConfig, error) {
	c := &httpServerConfig{}

	if cfg.CryptoKeyFile != "" {
		pKey, err := encryption.PrivateKeyFromFile(cfg.CryptoKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read crypro key: %w", err)
		}
		c.PKey = pKey
	} else {
		c.PKey = nil
	}

	if cfg.TrustedSubnet != "" {
		_, sn, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse trusted subnet: %w", err)
		}
		c.Subnet = sn
	} else {
		c.Subnet = nil
	}

	c.Key = cfg.Key
	c.Addr = cfg.RunAddr
	c.Storage = st

	return c, nil
}

func New(st MetricsStorage, cfg *config.ServerConfig) (*HTTPServer, error) {
	c, err := newConfig(st, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to configure http server: %w", err)
	}
	return newHTTPServer(newServer(configuration(c)...)), nil
}
