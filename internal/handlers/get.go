// Package handlers consist of handlers for the http server
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/storage"
)

// Get returns value of requested metric with text/plain Content-Type.
func Get(s MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ok, reason, code := validate(req, false)
		if !ok {
			http.Error(w, reason, code)
			return
		}

		metricType := chi.URLParam(req, "type")
		metricName := chi.URLParam(req, "name")

		switch metricType {
		case models.Counter:
			cvalue, err := s.GetCounter(req.Context(), metricName)
			if err != nil {

				status := http.StatusInternalServerError
				if errors.Is(err, storage.ErrMetricNotRegistered) {
					status = http.StatusNotFound
				}
				http.Error(w, err.Error(), status)
				return
			}
			_, err = w.Write([]byte(strconv.FormatInt(cvalue, 10)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case models.Gauge:
			gvalue, err := s.GetGauge(req.Context(), metricName)
			if err != nil {
				status := http.StatusInternalServerError
				if errors.Is(err, storage.ErrMetricNotRegistered) {
					status = http.StatusNotFound
				}
				http.Error(w, err.Error(), status)
				return
			}
			_, err = w.Write([]byte(strconv.FormatFloat(gvalue, 'f', -1, 64)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
