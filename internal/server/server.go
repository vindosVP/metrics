package server

import (
	"github.com/gorilla/mux"
	"github.com/vindosVP/metrics/internal/handlers"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"net/http"
)

func Run() error {

	gRepo := repos.NewGaugeRepo()
	cRepo := repos.NewCounterRepo()
	storage := memstorage.New(gRepo, cRepo)

	r := mux.NewRouter()
	r.HandleFunc("/update/{type}/{name}/{value}", handlers.Update(storage))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}

	return nil
}
