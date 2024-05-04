package loader

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/filestorage"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func BenchmarkLoader_LoadMetrics(b *testing.B) {

	ctx := context.Background()
	fName := "./bench-db.json"
	s := filestorage.NewFileStorage(repos.NewGaugeRepo(), repos.NewCounterRepo(), fName)
	defer os.Remove(fName)

	for i := 0; i < 100; i++ {
		s.UpdateCounter(ctx, RandStringRunes(10), rand.Int63())
		s.UpdateGauge(ctx, RandStringRunes(10), rand.Float64())
	}
	l := New(fName, s)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.LoadMetrics()
	}
}

func TestLoader(t *testing.T) {
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

	dump := &models.MetricsDump{}
	for k, v := range cMetrics {
		val := v
		metric := &models.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &val,
		}
		dump.Metrics = append(dump.Metrics, metric)
	}
	for k, v := range gMetrics {
		val := v
		metric := &models.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &val,
		}
		dump.Metrics = append(dump.Metrics, metric)
	}

	data, err := json.Marshal(dump)
	require.NoError(t, err)

	fileName := "./test_loader.json"
	err = os.WriteFile(fileName, data, 0666)
	defer os.Remove(fileName)
	require.NoError(t, err)

	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)
	loader := New(fileName, storage)
	err = loader.LoadMetrics()
	require.NoError(t, err)

	gotCMetrics, err := storage.GetAllCounter(context.Background())
	require.NoError(t, err)
	gotGMetrics, err := storage.GetAllGauge(context.Background())
	require.NoError(t, err)

	cEqual := reflect.DeepEqual(cMetrics, gotCMetrics)
	gEqual := reflect.DeepEqual(gMetrics, gotGMetrics)
	assert.True(t, cEqual)
	assert.True(t, gEqual)
}
