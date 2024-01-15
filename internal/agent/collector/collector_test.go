package collector

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test because testing.Short is enabled")
	}

	done := make(chan struct{})
	collectorShutdown := make(chan struct{})

	cfg := config.NewAgentConfig()
	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)
	c := New(cfg, storage)
	c.Done = done
	c.PollInterval = 1

	go func() {
		defer close(collectorShutdown)
		c.Run()
	}()
	time.Sleep(2 * time.Second)
	close(done)
	<-collectorShutdown

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

	cGauges, err := storage.GetAllGauge()
	require.NoError(t, err)
	cCounters, err := storage.GetAllCounter()
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
