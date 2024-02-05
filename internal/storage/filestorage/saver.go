package filestorage

import (
	"encoding/json"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"os"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	GetAllGauge() (map[string]float64, error)
	GetAllCounter() (map[string]int64, error)
}

type Saver struct {
	FileName      string
	StoreInterval time.Duration
	Done          <-chan struct{}
	Storage       MetricsStorage
}

func NewSaver(filename string, storeInterval time.Duration, s MetricsStorage) *Saver {
	return &Saver{
		FileName:      filename,
		StoreInterval: storeInterval,
		Storage:       s,
	}
}

func (s *Saver) Stop() {
	done := make(chan struct{})
	s.Done = done
	close(done)
}

func (s *Saver) Run() {
	tick := time.NewTicker(s.StoreInterval * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-s.Done:
			return
		case <-tick.C:
			s.Save()
		}
	}
}

func (s *Saver) Save() {
	gMetrics, err := s.Storage.GetAllGauge()
	if err != nil {
		logger.Log.Error("Failed to get gauge metrics", zap.Error(err))
	}
	cMetrics, err := s.Storage.GetAllCounter()
	if err != nil {
		logger.Log.Error("Failed to get counter metrics", zap.Error(err))
	}

	WriteMetrics(cMetrics, gMetrics, s.FileName)
}

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
