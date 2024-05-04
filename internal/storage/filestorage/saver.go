// Package filestorage is a package to save metrics values to inmemory storage and json file
package filestorage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
)

// MetricsStorage consists of methods to get metrics values from the storage
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	GetAllGauge(ctx context.Context) (map[string]float64, error)
	GetAllCounter(ctx context.Context) (map[string]int64, error)
}

// Saver consists data to save metrics dump
type Saver struct {
	Storage       MetricsStorage
	Done          <-chan struct{}
	FileName      string
	StoreInterval time.Duration
}

// NewSaver creates the Saver
func NewSaver(filename string, storeInterval time.Duration, s MetricsStorage) *Saver {
	return &Saver{
		FileName:      filename,
		StoreInterval: storeInterval,
		Storage:       s,
	}
}

// Stop method stops the Saver
func (s *Saver) Stop() {
	done := make(chan struct{})
	s.Done = done
	close(done)
}

// Run starts the Saver.
// Saver will get metrics from the storage and write them to the json file every n seconds.
func (s *Saver) Run() {
	tick := time.NewTicker(s.StoreInterval * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-s.Done:
			return
		case <-tick.C:
			s.save()
		}
	}
}

func (s *Saver) save() {
	ctx := context.Background()
	gMetrics, err := s.Storage.GetAllGauge(ctx)
	if err != nil {
		logger.Log.Error("Failed to get gauge metrics", zap.Error(err))
	}
	cMetrics, err := s.Storage.GetAllCounter(ctx)
	if err != nil {
		logger.Log.Error("Failed to get counter metrics", zap.Error(err))
	}

	WriteMetrics(cMetrics, gMetrics, s.FileName)
}

// WriteMetrics saves metrics values to file
func WriteMetrics(cMetrics map[string]int64, gMetrics map[string]float64, fileName string) {
	metrics := make([]*models.Metrics, len(gMetrics)+len(cMetrics))
	i := 0
	for k, v := range gMetrics {
		val := v
		metric := &models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &val,
		}
		metrics[i] = metric
		i++
	}
	for k, v := range cMetrics {
		val := v
		metric := &models.Metrics{
			ID:    k,
			MType: models.Counter,
			Delta: &val,
		}
		metrics[i] = metric
		i++
	}

	metricsDump := &models.MetricsDump{Metrics: metrics}

	data, err := json.MarshalIndent(metricsDump, "", "    ")
	if err != nil {
		logger.Log.Error("Failed to marshal metrics", zap.Error(err))
	}

	err = os.WriteFile(fileName, data, 0666)
	if err != nil {
		logger.Log.Error("Failed to write metrics", zap.Error(err))
	}
}
