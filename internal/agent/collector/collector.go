package collector

import (
	"context"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
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
	UpdateGauge(ctx context.Context, name string, v float64) (float64, error)
	UpdateCounter(ctx context.Context, name string, v int64) (int64, error)
	SetCounter(ctx context.Context, name string, v int64) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
	GetAllGauge(ctx context.Context) (map[string]float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAllCounter(ctx context.Context) (map[string]int64, error)
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
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.Log.Error("Failed to collect memory metrics", zap.Error(err))
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
