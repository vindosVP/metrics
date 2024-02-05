package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

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
		case models.Counter:
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
			logger.Log.Info("Updated metric value", zap.String("name", metricName), zap.Int64("value", cval))
		case models.Gauge:
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
			logger.Log.Info("Updated metric value", zap.String("name", metricName), zap.Float64("value", gval))
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
