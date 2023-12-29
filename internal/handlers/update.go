package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
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
		if req.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(req)
		metricType, ok := vars["type"]
		if !ok {
			http.Error(w, "type is missing in parameters", http.StatusBadRequest)
			return
		}

		if metricType != counter && metricType != gauge {
			http.Error(w, "invalid type parameter value", http.StatusBadRequest)
			return
		}

		metricName, ok := vars["name"]
		if !ok {
			http.Error(w, "name is missing in parameters", http.StatusNotFound)
			return
		}

		metricValue, ok := vars["value"]
		if !ok {
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
			log.Print(fmt.Sprintf("Updated value of %s with %v", metricName, cval))
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
			log.Print(fmt.Sprintf("Updated value of %s with %v", metricName, gval))
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
