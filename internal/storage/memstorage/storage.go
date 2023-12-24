package memstorage

import "github.com/vindosVP/metrics/internal/repos"

type Storage struct {
	gRepo *repos.GaugeRepo
	cRepo *repos.CounterRepo
}

func New(gRepo *repos.GaugeRepo, cRepo *repos.CounterRepo) *Storage {
	return &Storage{
		gRepo: gRepo,
		cRepo: cRepo,
	}
}

func (s Storage) UpdateGauge(name string, v float64) (float64, error) {
	return s.gRepo.Update(name, v)
}

func (s Storage) UpdateCounter(name string, v int64) (int64, error) {
	return s.cRepo.Update(name, v)
}

func (s Storage) GetGauge(name string) (float64, error) {
	return s.gRepo.Get(name)
}

func (s Storage) GetCounter(name string) (int64, error) {
	return s.cRepo.Get(name)
}
