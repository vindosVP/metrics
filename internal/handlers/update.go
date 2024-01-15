package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
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

func Update(s MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ok, reason, code := validate(req, true)
		if !ok {
			http.Error(w, reason, code)
			return
		}

		metricType := chi.URLParam(req, "type")
		metricName := chi.URLParam(req, "name")
		metricValue := chi.URLParam(req, "value")

		switch metricType {
		case counter:
			cval, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "invalid value type", http.StatusBadRequest)
				return
			}
			_, err = s.UpdateCounter(metricName, cval)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Updated value of %s with %v", metricName, cval)
		case gauge:
			gval, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "invalid value type", http.StatusBadRequest)
				return
			}
			_, err = s.UpdateGauge(metricName, gval)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Updated value of %s with %v", metricName, gval)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
