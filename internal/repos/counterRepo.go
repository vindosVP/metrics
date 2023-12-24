package repos

type CounterRepo struct {
	metrics map[string]int64
}

func NewCounterRepo() *CounterRepo {
	return &CounterRepo{metrics: make(map[string]int64)}
}

func (c CounterRepo) Update(name string, v int64) (int64, error) {
	currentV, ok := c.metrics[name]

	var newV int64
	if ok {
		newV = currentV + v
	} else {
		newV = v
	}
	c.metrics[name] = newV

	return newV, nil

}

func (c CounterRepo) Get(name string) (int64, error) {
	v, ok := c.metrics[name]
	if !ok {
		return 0, ErrMetricNotRegistered
	}

	return v, nil
}
