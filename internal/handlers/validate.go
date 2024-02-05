package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vindosVP/metrics/internal/models"
	"net/http"
)

func validate(req *http.Request, checkValue bool) (bool, string, int) {
	metricType := chi.URLParam(req, "type")
	if metricType == "" {
		return false, "type is missing in parameters", http.StatusBadRequest
	}
	if metricType != models.Counter && metricType != models.Gauge {
		return false, "invalid type parameter value", http.StatusBadRequest
	}
	metricName := chi.URLParam(req, "name")
	if metricName == "" {
		return false, "name is missing in parameters", http.StatusNotFound
	}
	if checkValue {
		metricValue := chi.URLParam(req, "value")
		if metricValue == "" {
			return false, "value is missing in parameters", http.StatusBadRequest
		}
	}

	return true, "", http.StatusOK
}

func validateUpdate(metrics *models.Metrics) (bool, string, int) {
	if metrics.MType != models.Counter && metrics.MType != models.Gauge {
		return false, "invalid metric type", http.StatusBadRequest
	}
	if metrics.ID == "" {
		return false, "invalid id", http.StatusNotFound
	}
	if metrics.MType == models.Counter && metrics.Delta == nil {
		return false, "invalid delta", http.StatusBadRequest
	}
	if metrics.MType == models.Gauge && metrics.Value == nil {
		return false, "invalid value", http.StatusBadRequest
	}
	return true, "", http.StatusOK
}

func validateGet(metrics *models.Metrics) (bool, string, int) {
	if metrics.MType != models.Counter && metrics.MType != models.Gauge {
		return false, "invalid metric type", http.StatusBadRequest
	}
	if metrics.ID == "" {
		return false, "invalid id", http.StatusNotFound
	}
	return true, "", http.StatusOK
}
