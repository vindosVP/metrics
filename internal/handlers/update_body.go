package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	UpdateGauge(name string, v float64) (float64, error)
	UpdateCounter(name string, v int64) (int64, error)
	SetCounter(name string, v int64) (int64, error)
	GetGauge(name string) (float64, error)
	GetAllGauge() (map[string]float64, error)
	GetCounter(name string) (int64, error)
	GetAllCounter() (map[string]int64, error)
}

func UpdateBody(s MetricsStorage) http.HandlerFunc {
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

		ok, reason, code := validateUpdate(metrics)
		if !ok {
			http.Error(w, reason, code)
			return
		}

		fields := []zap.Field{
			zap.String("name", metrics.ID),
			zap.String("type", metrics.MType),
		}
		resp := &models.Metrics{}

		switch metrics.MType {
		case models.Counter:
			delta := *metrics.Delta
			fields = append(fields, zap.Int64("delta", delta))
			val, err := s.UpdateCounter(metrics.ID, delta)
			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Log.Error("Failed to update metric value", fields...)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			resp.ID = metrics.ID
			resp.MType = models.Counter
			resp.Delta = &val

			logger.Log.Info("Updated metric value", fields...)
		case models.Gauge:
			value := *metrics.Value
			fields = append(fields, zap.Float64("value", value))
			val, err := s.UpdateGauge(metrics.ID, value)
			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Log.Error("Failed to update metric value", fields...)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			resp.ID = metrics.ID
			resp.MType = models.Gauge
			resp.Value = &val

			logger.Log.Info("Updated metric value", fields...)
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
