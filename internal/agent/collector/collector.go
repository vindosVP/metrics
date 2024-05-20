// Package collector collects metrics every n seconds.
package collector

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"

	"github.com/vindosVP/metrics/pkg/logger"
)

// Collector consists data to collect metrics.
type Collector struct {
	Done chan struct{}

	// Storage - storage to store metrics.
	Storage MetricsStorage

	// PollInterval - interval to collect metrics.
	PollInterval time.Duration
}

// MetricsStorage consists of methods to write and get data from storage.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {

	// UpdateGauge updates gauge metric.
	// new value replaces the old one.
	UpdateGauge(ctx context.Context, name string, v float64) (float64, error)

	// UpdateCounter updates counter metric.
	// new value adds to the old one.
	UpdateCounter(ctx context.Context, name string, v int64) (int64, error)
}

// New creates Collector.
func New(p time.Duration, s MetricsStorage) *Collector {
	return &Collector{
		Done:         make(chan struct{}),
		PollInterval: p,
		Storage:      s,
	}
}

func (c *Collector) Stop() {
	close(c.Done)
}

// Run starts the collector.
func (c *Collector) Run(wg *sync.WaitGroup) {
	tick := time.NewTicker(c.PollInterval * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-c.Done:
			wg.Done()
			return
		case <-tick.C:
			c.collectMetrics()
		}
	}
}

func (c *Collector) collectMetrics() {
	c.collectCounters()
	c.collectGauges()
	logger.Log.Info("Metrics collected")
}

func (c *Collector) collectCounters() {
	ctx := context.Background()
	_, err := c.Storage.UpdateCounter(ctx, "PollCount", 1)
	if err != nil {
		logger.Log.Error(
			"Failed to update metric",
			zap.String("name", "PollCount"),
			zap.Int64("value", 1),
			zap.Error(err))
	}
}

func (c *Collector) collectGauges() {
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	vMem, err := mem.VirtualMemory()
	if err != nil {
		logger.Log.Error("Failed to collect memory metrics", zap.Error(err))
	}
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.Log.Error("Failed to collect CPU metrics", zap.Error(err))
	}
	m := toMap(metrics, vMem, cpuUsage)
	ctx := context.Background()
	for key, val := range m {
		_, err := c.Storage.UpdateGauge(ctx, key, val)
		if err != nil {
			logger.Log.Error(
				"Failed to update metric",
				zap.String("name", key),
				zap.Float64("value", val),
				zap.Error(err))
		}
	}
}

func toMap(stats *runtime.MemStats, vMem *mem.VirtualMemoryStat, cpuUsage []float64) map[string]float64 {
	m := map[string]float64{
		"Alloc":           float64(stats.Alloc),
		"BuckHashSys":     float64(stats.BuckHashSys),
		"Frees":           float64(stats.Frees),
		"GCCPUFraction":   stats.GCCPUFraction,
		"GCSys":           float64(stats.GCSys),
		"HeapAlloc":       float64(stats.HeapAlloc),
		"HeapIdle":        float64(stats.HeapIdle),
		"HeapInuse":       float64(stats.HeapInuse),
		"HeapObjects":     float64(stats.HeapObjects),
		"HeapReleased":    float64(stats.HeapReleased),
		"HeapSys":         float64(stats.HeapSys),
		"LastGC":          float64(stats.LastGC),
		"Lookups":         float64(stats.Lookups),
		"MCacheSys":       float64(stats.MCacheSys),
		"MCacheInuse":     float64(stats.MCacheInuse),
		"NextGC":          float64(stats.NextGC),
		"MSpanInuse":      float64(stats.MSpanInuse),
		"MSpanSys":        float64(stats.MSpanSys),
		"Mallocs":         float64(stats.Mallocs),
		"NumForcedGC":     float64(stats.NumForcedGC),
		"NumGC":           float64(stats.NumGC),
		"OtherSys":        float64(stats.OtherSys),
		"PauseTotalNs":    float64(stats.PauseTotalNs),
		"StackInuse":      float64(stats.StackInuse),
		"StackSys":        float64(stats.StackSys),
		"Sys":             float64(stats.Sys),
		"TotalAlloc":      float64(stats.TotalAlloc),
		"RandomValue":     rand.Float64(),
		"TotalMemory":     float64(vMem.Total),
		"FreeMemory":      float64(vMem.Free),
		"CPUutilization1": cpuUsage[0],
	}

	return m
}
