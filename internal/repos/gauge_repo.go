package repos

import (
	"context"
	"sync"
)

// GaugeRepo - repository to store metrics with gauge type.
type GaugeRepo struct {
	metrics map[string]float64
	sync.Mutex
}

// NewGaugeRepo created the GaugeRepo
func NewGaugeRepo() *GaugeRepo {
	return &GaugeRepo{metrics: make(map[string]float64)}
}

// Update method updated value of gauge metric.
// new value replaces the old one.
func (g *GaugeRepo) Update(ctx context.Context, name string, v float64) (float64, error) {
	g.Lock()
	g.metrics[name] = v
	cVal := g.metrics[name]
	g.Unlock()
	return cVal, nil
}

// Get method returns gauge metric value
func (g *GaugeRepo) Get(ctx context.Context, name string) (float64, error) {
	g.Lock()
	v, ok := g.metrics[name]
	if !ok {
		g.Unlock()
		return 0, ErrMetricNotRegistered
	}
	g.Unlock()
	return v, nil
}

// GetAll method returns values of all collected gauge metrics
func (g *GaugeRepo) GetAll(ctx context.Context) (map[string]float64, error) {
	g.Lock()
	metrics := make(map[string]float64, len(g.metrics))
	for key, val := range g.metrics {
		metrics[key] = val
	}
	g.Unlock()
	return metrics, nil
}
