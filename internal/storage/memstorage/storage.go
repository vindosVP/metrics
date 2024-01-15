package memstorage

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

func (s *Storage) UpdateGauge(name string, v float64) (float64, error) {
	return s.gRepo.Update(name, v)
}

func (s *Storage) UpdateCounter(name string, v int64) (int64, error) {
	return s.cRepo.Update(name, v)
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
	return s.cRepo.Set(name, v)
}
