package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/storage"
	"github.com/vindosVP/metrics/pkg/logger"
)

// GetBody returns value of requested metric in json format with application/json Content-Type.
func GetBody(s MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		metrics := &models.Metrics{}
		var buf bytes.Buffer
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			logger.Log.Error("Failed to read request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			logger.Log.Error("Failed to unmarshal request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, reason, status := validateGet(metrics)
		if !ok {
			http.Error(w, reason, status)
			return
		}

		fields := []zap.Field{
			zap.String("name", metrics.ID),
			zap.String("type", metrics.MType),
		}
		resp := &models.Metrics{}

		switch metrics.MType {
		case models.Counter:
			val, cerr := s.GetCounter(req.Context(), metrics.ID)
			if cerr != nil {
				status = http.StatusInternalServerError
				if errors.Is(cerr, storage.ErrMetricNotRegistered) {
					status = http.StatusNotFound
				} else {
					fields = append(fields, zap.Error(cerr))
					logger.Log.Error("Failed to get metric value", fields...)
				}
				http.Error(w, cerr.Error(), status)
				return
			}

			resp.ID = metrics.ID
			resp.MType = models.Counter
			resp.Delta = &val
		case models.Gauge:
			val, gerr := s.GetGauge(req.Context(), metrics.ID)
			if gerr != nil {
				status = http.StatusInternalServerError
				if errors.Is(gerr, storage.ErrMetricNotRegistered) {
					status = http.StatusNotFound
				} else {
					fields = append(fields, zap.Error(gerr))
					logger.Log.Error("Failed to get metric value", fields...)
				}
				http.Error(w, gerr.Error(), status)
				return
			}

			resp.ID = metrics.ID
			resp.MType = models.Gauge
			resp.Value = &val
		}

		respData, err := json.Marshal(resp)
		if err != nil {
			logger.Log.Error("Failed to marshal response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(respData)
		if err != nil {
			logger.Log.Error("Failed to write response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
