package agent

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func Test_save(t *testing.T) {

	tests := []struct {
		name             string
		initialGauge     map[string]float64
		initialPollCount int64
	}{
		{
			name:             "free storage",
			initialGauge:     make(map[string]float64),
			initialPollCount: 0,
		},
		{
			name:             "initial poll count",
			initialGauge:     make(map[string]float64),
			initialPollCount: 12,
		},
		{
			name: "initial metrics",
			initialGauge: map[string]float64{
				"Alloc": 12.3,
			},
			initialPollCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cRepo := make(map[string]int64)
			cRepo["PollCount"] = tt.initialPollCount
			s := &storage{
				gaugeMetrics:   tt.initialGauge,
				counterMetrics: cRepo,
			}

			metrics := &runtime.MemStats{}
			runtime.ReadMemStats(metrics)
			save(s, metrics)

			assert.Equal(t, float64(metrics.Alloc), s.gaugeMetrics["Alloc"])
			assert.Equal(t, float64(metrics.BuckHashSys), s.gaugeMetrics["BuckHashSys"])
			assert.Equal(t, float64(metrics.Frees), s.gaugeMetrics["Frees"])
			assert.Equal(t, metrics.GCCPUFraction, s.gaugeMetrics["GCCPUFraction"])
			assert.Equal(t, float64(metrics.GCSys), s.gaugeMetrics["GCSys"])
			assert.Equal(t, float64(metrics.HeapAlloc), s.gaugeMetrics["HeapAlloc"])
			assert.Equal(t, float64(metrics.HeapIdle), s.gaugeMetrics["HeapIdle"])
			assert.Equal(t, float64(metrics.HeapInuse), s.gaugeMetrics["HeapInuse"])
			assert.Equal(t, float64(metrics.HeapObjects), s.gaugeMetrics["HeapObjects"])
			assert.Equal(t, float64(metrics.HeapReleased), s.gaugeMetrics["HeapReleased"])
			assert.Equal(t, float64(metrics.HeapSys), s.gaugeMetrics["HeapSys"])
			assert.Equal(t, float64(metrics.LastGC), s.gaugeMetrics["LastGC"])
			assert.Equal(t, float64(metrics.Lookups), s.gaugeMetrics["Lookups"])
			assert.Equal(t, float64(metrics.MCacheSys), s.gaugeMetrics["MCacheSys"])
			assert.Equal(t, float64(metrics.MSpanInuse), s.gaugeMetrics["MSpanInuse"])
			assert.Equal(t, float64(metrics.MSpanSys), s.gaugeMetrics["MSpanSys"])
			assert.Equal(t, float64(metrics.Mallocs), s.gaugeMetrics["Mallocs"])
			assert.Equal(t, float64(metrics.NumForcedGC), s.gaugeMetrics["NumForcedGC"])
			assert.Equal(t, float64(metrics.NumGC), s.gaugeMetrics["NumGC"])
			assert.Equal(t, float64(metrics.OtherSys), s.gaugeMetrics["OtherSys"])
			assert.Equal(t, float64(metrics.PauseTotalNs), s.gaugeMetrics["PauseTotalNs"])
			assert.Equal(t, float64(metrics.StackInuse), s.gaugeMetrics["StackInuse"])
			assert.Equal(t, float64(metrics.StackSys), s.gaugeMetrics["StackSys"])
			assert.Equal(t, float64(metrics.Sys), s.gaugeMetrics["Sys"])
			assert.Equal(t, float64(metrics.TotalAlloc), s.gaugeMetrics["TotalAlloc"])

			_, ok := s.gaugeMetrics["RandomValue"]
			assert.True(t, ok)
			assert.Equal(t, tt.initialPollCount+1, s.counterMetrics["PollCount"])

		})
	}

}
