package agent

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func collect(s *storage, wg *sync.WaitGroup) {
	for {
		s.mu.Lock()

		metrics := &runtime.MemStats{}
		runtime.ReadMemStats(metrics)
		save(s, metrics)
		log.Print("Metrics collectd")

		s.mu.Unlock()

		time.Sleep(pollInterval * time.Second)
	}
}

func save(s *storage, metrics *runtime.MemStats) {
	s.gaugeMetrics["Alloc"] = float64(metrics.Alloc)
	s.gaugeMetrics["BuckHashSys"] = float64(metrics.BuckHashSys)
	s.gaugeMetrics["Frees"] = float64(metrics.Frees)
	s.gaugeMetrics["GCCPUFraction"] = metrics.GCCPUFraction
	s.gaugeMetrics["GCSys"] = float64(metrics.GCSys)
	s.gaugeMetrics["HeapAlloc"] = float64(metrics.HeapAlloc)
	s.gaugeMetrics["HeapIdle"] = float64(metrics.HeapIdle)
	s.gaugeMetrics["HeapInuse"] = float64(metrics.HeapInuse)
	s.gaugeMetrics["HeapObjects"] = float64(metrics.HeapObjects)
	s.gaugeMetrics["HeapReleased"] = float64(metrics.HeapReleased)
	s.gaugeMetrics["HeapSys"] = float64(metrics.HeapSys)
	s.gaugeMetrics["LastGC"] = float64(metrics.LastGC)
	s.gaugeMetrics["Lookups"] = float64(metrics.Lookups)
	s.gaugeMetrics["MCacheSys"] = float64(metrics.MCacheSys)
	s.gaugeMetrics["MSpanInuse"] = float64(metrics.MSpanInuse)
	s.gaugeMetrics["MSpanSys"] = float64(metrics.MSpanSys)
	s.gaugeMetrics["Mallocs"] = float64(metrics.Mallocs)
	s.gaugeMetrics["NumForcedGC"] = float64(metrics.NumForcedGC)
	s.gaugeMetrics["NumGC"] = float64(metrics.NumGC)
	s.gaugeMetrics["OtherSys"] = float64(metrics.OtherSys)
	s.gaugeMetrics["PauseTotalNs"] = float64(metrics.PauseTotalNs)
	s.gaugeMetrics["StackInuse"] = float64(metrics.StackInuse)
	s.gaugeMetrics["StackSys"] = float64(metrics.StackSys)
	s.gaugeMetrics["Sys"] = float64(metrics.Sys)
	s.gaugeMetrics["TotalAlloc"] = float64(metrics.TotalAlloc)
	s.gaugeMetrics["RandomValue"] = rand.Float64()
	s.counterMetrics["PollCount"] += 1
}
