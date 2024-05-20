package collector

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
)

func BenchmarkCollectorCollectMetrics(b *testing.B) {
	c := New(10*time.Second, memstorage.New(repos.NewGaugeRepo(), repos.NewCounterRepo()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.collectMetrics()
	}
}

func TestCollector(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test because testing.Short is enabled")
	}

	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)
	c := New(10*time.Second, storage)
	c.PollInterval = 1

	wg := sync.WaitGroup{}
	wg.Add(1)
	go c.Run(&wg)
	time.Sleep(2 * time.Second)
	c.Stop()
	wg.Wait()

	expGauges := []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
	}

	expCounters := []string{"PollCount"}

	ctx := context.Background()
	cGauges, err := storage.GetAllGauge(ctx)
	require.NoError(t, err)
	cCounters, err := storage.GetAllCounter(ctx)
	require.NoError(t, err)

	for _, key := range expGauges {
		_, ok := cGauges[key]
		assert.True(t, ok)
	}
	for _, key := range expCounters {
		_, ok := cCounters[key]
		assert.True(t, ok)
	}
}
