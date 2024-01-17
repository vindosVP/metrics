package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config2 "github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/handlers"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"log"
	"net/http"
)

func Run(cfg *config2.ServerConfig) error {

	gRepo := repos.NewGaugeRepo()
	cRepo := repos.NewCounterRepo()
	storage := memstorage.New(gRepo, cRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{type}/{name}/{value}", handlers.Update(storage))
	r.Get("/value/{type}/{name}", handlers.Get(storage))
	r.Get("/", handlers.List(storage))
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	log.Printf("Running server on %s", cfg.RunAddr)
	err := http.ListenAndServe(cfg.RunAddr, r)
	if err != nil {
		return err
	}

	return nil
}
