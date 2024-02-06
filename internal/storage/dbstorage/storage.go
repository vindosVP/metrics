package dbstorage

import (
	"context"
)

type Counter interface {
	Insert(ctx context.Context, name string, v int64) (int64, error)
	Update(ctx context.Context, name string, v int64) (int64, error)
	Exists(ctx context.Context, name string) (bool, error)
	Get(ctx context.Context, name string) (int64, error)
	GetAll(ctx context.Context) (map[string]int64, error)
}

type Gauge interface {
	Insert(ctx context.Context, name string, v float64) (float64, error)
	Update(ctx context.Context, name string, v float64) (float64, error)
	Exists(ctx context.Context, name string) (bool, error)
	Get(ctx context.Context, name string) (float64, error)
	GetAll(ctx context.Context) (map[string]float64, error)
}

type Storage struct {
	cr Counter
	gr Gauge
}

func New(cRepo Counter, gRepo Gauge) *Storage {
	return &Storage{
		cr: cRepo,
		gr: gRepo,
	}
}

func (s *Storage) UpdateGauge(ctx context.Context, name string, v float64) (float64, error) {
	exists, err := s.gr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if exists {
		return s.gr.Update(ctx, name, v)
	}
	return s.gr.Insert(ctx, name, v)
}

func (s *Storage) UpdateCounter(ctx context.Context, name string, v int64) (int64, error) {
	exists, err := s.cr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return s.cr.Insert(ctx, name, v)
	}
	cval, err := s.cr.Get(ctx, name)
	if err != nil {
		return 0, err
	}
	return s.cr.Update(ctx, name, cval+v)
}

func (s *Storage) GetGauge(ctx context.Context, name string) (float64, error) {
	exists, err := s.gr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrMetricNotRegistered
	}
	return s.gr.Get(ctx, name)
}

func (s *Storage) GetCounter(ctx context.Context, name string) (int64, error) {
	exists, err := s.cr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrMetricNotRegistered
	}
	return s.cr.Get(ctx, name)
}

func (s *Storage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return s.gr.GetAll(ctx)
}

func (s *Storage) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	return s.cr.GetAll(ctx)
}

func (s *Storage) SetCounter(ctx context.Context, name string, v int64) (int64, error) {
	exists, err := s.cr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return s.cr.Insert(ctx, name, v)
	}
	return s.cr.Update(ctx, name, v)
}
