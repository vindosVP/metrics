package collector

import (
	"github.com/vindosVP/metrics/internal/config"
	"github.com/vindosVP/metrics/internal/storage"
	"log"
	"math/rand"
	"runtime"
	"time"
)

type Collector struct {
	PollInterval int
	Done         <-chan struct{}
	Storage      storage.MetricsStorage
}

func New(cfg *config.AgentConfig, s storage.MetricsStorage) *Collector {
	return &Collector{
		PollInterval: cfg.PollInterval,
		Storage:      s,
	}
}

func (c *Collector) Run() {
	tick := time.NewTicker(time.Duration(c.PollInterval) * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-c.Done:
			return
		case <-tick.C:
			c.CollectMetrics()
		}
	}
}

func (c *Collector) CollectMetrics() {
	c.collectCounters()
	c.collectGauges()
	log.Print("Metrics collected")
}

func (c *Collector) collectCounters() {
	_, err := c.Storage.UpdateCounter("PollCount", 1)
	if err != nil {
		log.Printf("Failed to update counter PollCount: %v", err)
	}
}

func (c *Collector) collectGauges() {
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	m := toMap(metrics)
	for key, val := range m {
		_, err := c.Storage.UpdateGauge(key, val)
		if err != nil {
			log.Printf("Failed to gauge %s: %v", key, err)
		}
	}
}

func toMap(stats *runtime.MemStats) map[string]float64 {
	m := make(map[string]float64, 26)
	m["Alloc"] = float64(stats.Alloc)
	m["BuckHashSys"] = float64(stats.BuckHashSys)
	m["Frees"] = float64(stats.Frees)
	m["GCCPUFraction"] = stats.GCCPUFraction
	m["GCSys"] = float64(stats.GCSys)
	m["HeapAlloc"] = float64(stats.HeapAlloc)
	m["HeapIdle"] = float64(stats.HeapIdle)
	m["HeapInuse"] = float64(stats.HeapInuse)
	m["HeapObjects"] = float64(stats.HeapObjects)
	m["HeapReleased"] = float64(stats.HeapReleased)
	m["HeapSys"] = float64(stats.HeapSys)
	m["LastGC"] = float64(stats.LastGC)
	m["Lookups"] = float64(stats.Lookups)
	m["MCacheSys"] = float64(stats.MCacheSys)
	m["MSpanInuse"] = float64(stats.MSpanInuse)
	m["MSpanSys"] = float64(stats.MSpanSys)
	m["Mallocs"] = float64(stats.Mallocs)
	m["NumForcedGC"] = float64(stats.NumForcedGC)
	m["NumGC"] = float64(stats.NumGC)
	m["OtherSys"] = float64(stats.OtherSys)
	m["PauseTotalNs"] = float64(stats.PauseTotalNs)
	m["StackInuse"] = float64(stats.StackInuse)
	m["StackSys"] = float64(stats.StackSys)
	m["Sys"] = float64(stats.Sys)
	m["TotalAlloc"] = float64(stats.TotalAlloc)
	m["RandomValue"] = rand.Float64()
	return m
}
