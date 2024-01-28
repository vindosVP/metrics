package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/handlers"
	"github.com/vindosVP/metrics/internal/middleware"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/server/loader"
	"github.com/vindosVP/metrics/internal/server/saver"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"log"
	"net/http"
)

func Run(cfg *config.ServerConfig) error {

	gRepo := repos.NewGaugeRepo()
	cRepo := repos.NewCounterRepo()
	storage := memstorage.New(gRepo, cRepo)

	if cfg.Restore {
		dumpLoader := loader.New(cfg.FileStoragePath, storage)
		err := dumpLoader.LoadMetrics()
		if err != nil {
			return err
		}
	}

	r := chi.NewRouter()
	r.Post("/update/", middleware.WithLogging(middleware.WithCompression(handlers.UpdateBody(storage))))
	r.Post("/value/", middleware.WithLogging(middleware.WithCompression(handlers.GetBody(storage))))
	r.Post("/update/{type}/{name}/{value}", middleware.WithLogging(handlers.Update(storage)))
	r.Get("/value/{type}/{name}", middleware.WithLogging(handlers.Get(storage)))
	r.Get("/", middleware.WithLogging(middleware.WithCompression(handlers.List(storage))))
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	svr := saver.New(cfg.FileStoragePath, cfg.StoreInterval, storage)
	go svr.Run()
	defer svr.Stop()

	log.Printf("Running server on %s", cfg.RunAddr)
	err := http.ListenAndServe(cfg.RunAddr, r)
	if err != nil {
		return err
	}

	return nil
}
