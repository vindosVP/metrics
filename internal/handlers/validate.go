package handlers

import (
	"github.com/vindosVP/metrics/internal/models"
	"net/http"
)

const (
	counter = "counter"
	gauge   = "gauge"
)

func validateUpdate(metrics *models.Metrics) (bool, string, int) {
	if metrics.MType != counter && metrics.MType != gauge {
		return false, "invalid metric type", http.StatusBadRequest
	}
	if metrics.ID == "" {
		return false, "invalid id", http.StatusNotFound
	}
	if metrics.MType == counter && metrics.Delta == nil {
		return false, "invalid delta", http.StatusBadRequest
	}
	if metrics.MType == gauge && metrics.Value == nil {
		return false, "invalid value", http.StatusBadRequest
	}
	return true, "", http.StatusOK
}

func validateGet(metrics *models.Metrics) (bool, string, int) {
	if metrics.MType != counter && metrics.MType != gauge {
		return false, "invalid metric type", http.StatusBadRequest
	}
	if metrics.ID == "" {
		return false, "invalid id", http.StatusNotFound
	}
	return true, "", http.StatusOK
}
