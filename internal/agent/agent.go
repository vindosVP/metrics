package agent

import (
	"sync"
)

const (
	gauge   = metricType("gauge")
	counter = metricType("counter")
)

type metricType string

type storage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	mu             sync.Mutex
}

func Run() error {

	cfg := NewConfig()

	s := &storage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
	s.counterMetrics["PollCount"] = 0

	var wg sync.WaitGroup

	wg.Add(2)
	go collect(s, cfg, &wg)
	go send(s, cfg, &wg)
	wg.Wait()

	return nil
}
