package repos

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Counter
type Counter interface {
	Update(name string, v int64) (int64, error)
	Get(name string) (int64, error)
	GetAll() (map[string]int64, error)
	Set(name string, v int64) (int64, error)
}

type CounterRepo struct {
	metrics map[string]int64
}

func NewCounterRepo() *CounterRepo {
	return &CounterRepo{metrics: make(map[string]int64)}
}

func (c *CounterRepo) Update(name string, v int64) (int64, error) {
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

func (c *CounterRepo) Get(name string) (int64, error) {
	v, ok := c.metrics[name]
	if !ok {
		return 0, ErrMetricNotRegistered
	}

	return v, nil
}

func (c *CounterRepo) GetAll() (map[string]int64, error) {
	return c.metrics, nil
}

func (c *CounterRepo) Set(name string, v int64) (int64, error) {
	c.metrics[name] = v
	return c.metrics[name], nil
}
