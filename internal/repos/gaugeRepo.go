package repos

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Gauge
type Gauge interface {
	Update(name string, v float64) (float64, error)
	Get(name string) (float64, error)
	GetAll() (map[string]float64, error)
}

type GaugeRepo struct {
	metrics map[string]float64
}

func NewGaugeRepo() *GaugeRepo {
	return &GaugeRepo{metrics: make(map[string]float64)}
}

func (g *GaugeRepo) Update(name string, v float64) (float64, error) {
	g.metrics[name] = v
	return v, nil
}

func (g *GaugeRepo) Get(name string) (float64, error) {
	v, ok := g.metrics[name]
	if !ok {
		return 0, ErrMetricNotRegistered
	}

	return v, nil
}

func (g *GaugeRepo) GetAll() (map[string]float64, error) {
	return g.metrics, nil
}
