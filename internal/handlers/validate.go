package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	counter = "counter"
	gauge   = "gauge"
)

func validate(req *http.Request, checkValue bool) (bool, string, int) {

	metricType := chi.URLParam(req, "type")
	if metricType == "" {
		return false, "type is missing in parameters", http.StatusBadRequest
	}

	if metricType != counter && metricType != gauge {
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
