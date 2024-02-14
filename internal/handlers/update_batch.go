package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"net/http"
)

func UpdateBatch(s MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		batch := make([]*models.Metrics, 0)
		var buf bytes.Buffer
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			logger.Log.Error("Failed to read request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &batch); err != nil {
			logger.Log.Error("Failed to unmarshal request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i, metric := range batch {
			ok, reason, code := validateUpdate(metric)
			if !ok {
				http.Error(w, fmt.Sprintf("bad structure number %d: %s", i, reason), code)
				return
			}
		}

		err = s.InsertBatch(req.Context(), batch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
