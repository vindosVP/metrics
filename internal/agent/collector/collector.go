package collector

import (
	"github.com/vindosVP/metrics/cmd/agent/config"
	"log"
	"math/rand"
	"runtime"
	"time"
)

type Collector struct {
	PollInterval time.Duration
	Done         <-chan struct{}
	Storage      MetricsStorage
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	UpdateGauge(name string, v float64) (float64, error)
	UpdateCounter(name string, v int64) (int64, error)
	SetCounter(name string, v int64) (int64, error)
	GetGauge(name string) (float64, error)
	GetAllGauge() (map[string]float64, error)
	GetCounter(name string) (int64, error)
	GetAllCounter() (map[string]int64, error)
}

func New(cfg *config.AgentConfig, s MetricsStorage) *Collector {
	return &Collector{
		PollInterval: cfg.PollInterval,
		Storage:      s,
	}
}

func (c *Collector) Run() {
	tick := time.NewTicker(c.PollInterval * time.Second)
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
	m := map[string]float64{
		"Alloc":         float64(stats.Alloc),
		"BuckHashSys":   float64(stats.BuckHashSys),
		"Frees":         float64(stats.Frees),
		"GCCPUFraction": stats.GCCPUFraction,
		"GCSys":         float64(stats.GCSys),
		"HeapAlloc":     float64(stats.HeapAlloc),
		"HeapIdle":      float64(stats.HeapIdle),
		"HeapInuse":     float64(stats.HeapInuse),
		"HeapObjects":   float64(stats.HeapObjects),
		"HeapReleased":  float64(stats.HeapReleased),
		"HeapSys":       float64(stats.HeapSys),
		"LastGC":        float64(stats.LastGC),
		"Lookups":       float64(stats.Lookups),
		"MCacheSys":     float64(stats.MCacheSys),
		"MSpanInuse":    float64(stats.MSpanInuse),
		"MSpanSys":      float64(stats.MSpanSys),
		"Mallocs":       float64(stats.Mallocs),
		"NumForcedGC":   float64(stats.NumForcedGC),
		"NumGC":         float64(stats.NumGC),
		"OtherSys":      float64(stats.OtherSys),
		"PauseTotalNs":  float64(stats.PauseTotalNs),
		"StackInuse":    float64(stats.StackInuse),
		"StackSys":      float64(stats.StackSys),
		"Sys":           float64(stats.Sys),
		"TotalAlloc":    float64(stats.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}

	return m
}
