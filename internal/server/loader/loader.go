package loader

import (
	"context"
	"encoding/json"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"os"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	SetCounter(ctx context.Context, name string, v int64) (int64, error)
	UpdateGauge(ctx context.Context, name string, v float64) (float64, error)
}

type Loader struct {
	filename string
	storage  MetricsStorage
}

func New(filename string, storage MetricsStorage) *Loader {
	return &Loader{
		filename: filename,
		storage:  storage,
	}
}

func (l *Loader) LoadMetrics() error {
	ctx := context.Background()
	data, err := os.ReadFile(l.filename)
	if err != nil {
		return err
	}
	metricsDump := &models.MetricsDump{}
	err = json.Unmarshal(data, &metricsDump)
	if err != nil {
		return err
	}
	for _, metric := range metricsDump.Metrics {
		if metric.MType == models.Counter {
			_, err := l.storage.SetCounter(ctx, metric.ID, *metric.Delta)
			if err != nil {
				logger.Log.Error("Failed to set counter", zap.Error(err))
			}
		} else if metric.MType == models.Gauge {
			_, err := l.storage.UpdateGauge(ctx, metric.ID, *metric.Value)
			if err != nil {
				logger.Log.Error("Failed to update counter", zap.Error(err))
			}
		}
	}
	return nil
}
