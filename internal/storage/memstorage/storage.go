package memstorage

import (
	"context"
	"errors"

	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage"
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

type Storage struct {
	gRepo Gauge
	cRepo Counter
}

func New(gRepo Gauge, cRepo Counter) *Storage {
	return &Storage{
		gRepo: gRepo,
		cRepo: cRepo,
	}
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
	return nil
}

func (s *Storage) UpdateGauge(ctx context.Context, name string, v float64) (float64, error) {
	return s.gRepo.Update(ctx, name, v)
}

func (s *Storage) UpdateCounter(ctx context.Context, name string, v int64) (int64, error) {
	return s.cRepo.Update(ctx, name, v)
}

func (s *Storage) GetGauge(ctx context.Context, name string) (float64, error) {
	val, err := s.gRepo.Get(ctx, name)
	if errors.Is(err, repos.ErrMetricNotRegistered) {
		return 0, storage.ErrMetricNotRegistered
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *Storage) GetCounter(ctx context.Context, name string) (int64, error) {
	val, err := s.cRepo.Get(ctx, name)
	if errors.Is(err, repos.ErrMetricNotRegistered) {
		return 0, storage.ErrMetricNotRegistered
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *Storage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return s.gRepo.GetAll(ctx)
}

func (s *Storage) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	return s.cRepo.GetAll(ctx)
}

func (s *Storage) SetCounter(ctx context.Context, name string, v int64) (int64, error) {
	return s.cRepo.Set(ctx, name, v)
}
