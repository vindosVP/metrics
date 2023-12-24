package repos

type GaugeRepo struct {
	metrics map[string]float64
}

func NewGaugeRepo() *GaugeRepo {
	return &GaugeRepo{metrics: make(map[string]float64)}
}

func (g GaugeRepo) Update(name string, v float64) (float64, error) {
	currentV, ok := g.metrics[name]

	var newV float64
	if ok {
		newV = currentV + v
	} else {
		newV = v
	}
	g.metrics[name] = newV

	return newV, nil
}

func (g GaugeRepo) Get(name string) (float64, error) {
	v, ok := g.metrics[name]
	if !ok {
		return 0, ErrMetricNotRegistered
	}

	return v, nil
}
