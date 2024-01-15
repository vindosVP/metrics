package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vindosVP/metrics/internal/repos"
	"net/http"
	"strconv"
)

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
