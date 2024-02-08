package filestorage

import (
	"context"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Counter
type Counter interface {
	Update(ctx context.Context, name string, v int64) (int64, error)
	Get(ctx context.Context, name string) (int64, error)
	GetAll(ctx context.Context) (map[string]int64, error)
	Set(ctx context.Context, name string, v int64) (int64, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Gauge
type Gauge interface {
	Update(ctx context.Context, name string, v float64) (float64, error)
	Get(ctx context.Context, name string) (float64, error)
	GetAll(ctx context.Context) (map[string]float64, error)
}

func NewFileStorage(gRepo Gauge, cRepo Counter, fileName string) *Storage {
	return &Storage{
		gRepo:    gRepo,
		cRepo:    cRepo,
		fileName: fileName,
	}
}

type Storage struct {
	gRepo    Gauge
	cRepo    Counter
	fileName string
}

func (s *Storage) InsertBatch(ctx context.Context, batch []*models.Metrics) error {
	for _, metric := range batch {
		switch metric.MType {
		case models.Counter:
			val := *metric.Delta
			_, err := s.cRepo.Update(ctx, metric.ID, val)
			if err != nil {
				return err
			}
		case models.Gauge:
			val := *metric.Value
			_, err := s.gRepo.Update(ctx, metric.ID, val)
			if err != nil {
				return err
			}
		}
	}
	s.dump(ctx)
	return nil
}

func (s *Storage) UpdateGauge(ctx context.Context, name string, v float64) (float64, error) {
	val, err := s.gRepo.Update(ctx, name, v)
	s.dump(ctx)
	return val, err
}

func (s *Storage) UpdateCounter(ctx context.Context, name string, v int64) (int64, error) {
	val, err := s.cRepo.Update(ctx, name, v)
	s.dump(ctx)
	return val, err
}

func (s *Storage) GetGauge(ctx context.Context, name string) (float64, error) {
	return s.gRepo.Get(ctx, name)
}

func (s *Storage) GetCounter(ctx context.Context, name string) (int64, error) {
	return s.cRepo.Get(ctx, name)
}

func (s *Storage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return s.gRepo.GetAll(ctx)
}

func (s *Storage) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	return s.cRepo.GetAll(ctx)
}

func (s *Storage) SetCounter(ctx context.Context, name string, v int64) (int64, error) {
	val, err := s.cRepo.Set(ctx, name, v)
	s.dump(ctx)
	return val, err
}

func (s *Storage) dump(ctx context.Context) {
	cMetrics, err := s.GetAllCounter(ctx)
	if err != nil {
		logger.Log.Error("Failed to get counters", zap.Error(err))
	}
	gMetrics, err := s.GetAllGauge(ctx)
	if err != nil {
		logger.Log.Error("Failed to get gauges", zap.Error(err))
	}
	WriteMetrics(cMetrics, gMetrics, s.fileName)
}
