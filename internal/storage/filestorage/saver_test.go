package filestorage

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
)

func TestSaver(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test because testing.Short is enabled")
	}

	cMetrics := map[string]int64{
		"Test1": 10,
		"Test2": 15,
		"Test3": 20,
	}
	gMetrics := map[string]float64{
		"Test1": 1.5,
		"Test2": 3,
		"Test3": 1112.45,
	}
	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)

	ctx := context.Background()
	for k, v := range cMetrics {
		_, err := storage.SetCounter(ctx, k, v)
		require.NoError(t, err)
	}
	for k, v := range gMetrics {
		_, err := storage.UpdateGauge(ctx, k, v)
		require.NoError(t, err)
	}

	fileName := "./test_saver.json"
	interval := time.Duration(1)
	saver := NewSaver(fileName, interval, storage)
	defer os.Remove(fileName)

	done := make(chan struct{})
	saverShutdown := make(chan struct{})

	saver.Done = done

	go func() {
		defer close(saverShutdown)
		saver.Run()
	}()
	time.Sleep(2 * time.Second)
	close(done)
	<-saverShutdown

	savedData, err := os.ReadFile(fileName)
	require.NoError(t, err)
	dump := &models.MetricsDump{}
	err = json.Unmarshal(savedData, &dump)
	require.NoError(t, err)

	gotCMetrics := make(map[string]int64)
	gotGMetrics := make(map[string]float64)

	for _, metric := range dump.Metrics {
		if metric.MType == "counter" {
			gotCMetrics[metric.ID] = *metric.Delta
		} else {
			gotGMetrics[metric.ID] = *metric.Value
		}
	}

	cEqual := reflect.DeepEqual(cMetrics, gotCMetrics)
	gEqual := reflect.DeepEqual(gMetrics, gotGMetrics)
	assert.True(t, cEqual)
	assert.True(t, gEqual)
}
