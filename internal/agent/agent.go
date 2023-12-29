package agent

import (
	"sync"
)

const (
	pollInterval   = 2
	reportInterval = 10
)

const (
	gauge   = metricType("gauge")
	counter = metricType("counter")
)

const serverHost = "localhost:8080"

type metricType string

type storage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	mu             sync.Mutex
}

func Run() error {

	s := &storage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
	s.counterMetrics["PollCount"] = 0

	var wg sync.WaitGroup

	wg.Add(2)
	go collect(s, &wg)
	go send(s, &wg)
	wg.Wait()

	return nil
}
