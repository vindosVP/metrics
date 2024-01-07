package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vindosVP/metrics/internal/storage"
	"log"
	"net/http"
	"strconv"
)

const (
	counter = "counter"
	gauge   = "gauge"
)

func Update(s storage.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		metricType := chi.URLParam(req, "type")
		if metricType == "" {
			http.Error(w, "type is missing in parameters", http.StatusBadRequest)
			return
		}

		if metricType != counter && metricType != gauge {
			http.Error(w, "invalid type parameter value", http.StatusBadRequest)
			return
		}

		metricName := chi.URLParam(req, "name")
		if metricName == "" {
			http.Error(w, "name is missing in parameters", http.StatusNotFound)
			return
		}

		metricValue := chi.URLParam(req, "value")
		if metricValue == "" {
			http.Error(w, "value is missing in parameters", http.StatusBadRequest)
			return
		}

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
