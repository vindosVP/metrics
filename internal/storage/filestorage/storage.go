package filestorage

import (
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Counter
type Counter interface {
	Update(name string, v int64) (int64, error)
	Get(name string) (int64, error)
	GetAll() (map[string]int64, error)
	Set(name string, v int64) (int64, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Gauge
type Gauge interface {
	Update(name string, v float64) (float64, error)
	Get(name string) (float64, error)
	GetAll() (map[string]float64, error)
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

func (s *Storage) UpdateGauge(name string, v float64) (float64, error) {
	val, err := s.gRepo.Update(name, v)
	s.dump()
	return val, err
}

func (s *Storage) UpdateCounter(name string, v int64) (int64, error) {
	val, err := s.cRepo.Update(name, v)
	s.dump()
	return val, err
}

func (s *Storage) GetGauge(name string) (float64, error) {
	return s.gRepo.Get(name)
}

func (s *Storage) GetCounter(name string) (int64, error) {
	return s.cRepo.Get(name)
}

func (s *Storage) GetAllGauge() (map[string]float64, error) {
	return s.gRepo.GetAll()
}

func (s *Storage) GetAllCounter() (map[string]int64, error) {
	return s.cRepo.GetAll()
}

func (s *Storage) SetCounter(name string, v int64) (int64, error) {
	val, err := s.cRepo.Set(name, v)
	s.dump()
	return val, err
}

func (s *Storage) dump() {
	cMetrics, err := s.GetAllCounter()
	if err != nil {
		logger.Log.Error("Failed to get counters", zap.Error(err))
	}
	gMetrics, err := s.GetAllGauge()
	if err != nil {
		logger.Log.Error("Failed to get gauges", zap.Error(err))
	}
	WriteMetrics(cMetrics, gMetrics, s.fileName)
}
