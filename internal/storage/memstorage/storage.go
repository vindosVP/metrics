package memstorage

import "context"

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

func (s *Storage) UpdateGauge(ctx context.Context, name string, v float64) (float64, error) {
	return s.gRepo.Update(ctx, name, v)
}

func (s *Storage) UpdateCounter(ctx context.Context, name string, v int64) (int64, error) {
	return s.cRepo.Update(ctx, name, v)
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
	return s.cRepo.Set(ctx, name, v)
}
