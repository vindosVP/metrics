package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage"
	"net/http"
	"strconv"
)

func Get(s storage.MetricsStorage) http.HandlerFunc {
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

		switch metricType {
		case counter:
			cvalue, err := s.GetCounter(metricName)
			if err != nil {
				var status int
				if err == repos.ErrMetricNotRegistered {
					status = http.StatusNotFound
				} else {
					status = http.StatusInternalServerError
				}
				http.Error(w, err.Error(), status)
				return
			}
			_, err = w.Write([]byte(strconv.FormatInt(cvalue, 10)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case gauge:
			gvalue, err := s.GetGauge(metricName)
			if err != nil {
				var status int
				if err == repos.ErrMetricNotRegistered {
					status = http.StatusNotFound
				} else {
					status = http.StatusInternalServerError
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
