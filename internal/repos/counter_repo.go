package repos

import "sync"

type CounterRepo struct {
	metrics map[string]int64
	sync.Mutex
}

func NewCounterRepo() *CounterRepo {
	return &CounterRepo{metrics: make(map[string]int64)}
}

func (c *CounterRepo) Update(name string, v int64) (int64, error) {
	c.Lock()
	currentV, ok := c.metrics[name]

	var newV int64
	if ok {
		newV = currentV + v
	} else {
		newV = v
	}
	c.metrics[name] = newV
	c.Unlock()
	return newV, nil

}

func (c *CounterRepo) Get(name string) (int64, error) {
	c.Lock()
	v, ok := c.metrics[name]
	if !ok {
		c.Unlock()
		return 0, ErrMetricNotRegistered
	}
	c.Unlock()
	return v, nil
}

func (c *CounterRepo) GetAll() (map[string]int64, error) {
	c.Lock()
	metrics := make(map[string]int64, len(c.metrics))
	for key, val := range c.metrics {
		metrics[key] = val
	}
	c.Unlock()
	return metrics, nil
}

func (c *CounterRepo) Set(name string, v int64) (int64, error) {
	c.Lock()
	c.metrics[name] = v
	c.Unlock()
	return c.metrics[name], nil
}
