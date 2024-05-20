package sender

import (
	"context"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func BenchmarkMakeButch(b *testing.B) {
	c := make(map[string]int64, 100)
	g := make(map[string]float64, 100)

	for i := 0; i < 100; i++ {
		c[RandStringRunes(10)] = rand.Int63()
		g[RandStringRunes(10)] = rand.Float64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeButch(c, g)
	}
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestSender(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test because testing.Short is enabled")
	}

	cfg := config.NewAgentConfig()
	cRepo := repos.NewCounterRepo()
	ctx := context.Background()
	_, err := cRepo.Set(ctx, "PollCount", 100)
	require.NoError(t, err)
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)
	c := New(cfg, storage, nil)
	c.ReportInterval = 1

	responder := httpmock.NewStringResponder(200, "")
	httpmock.RegisterResponder(http.MethodPost, `=~^(http|https)://.+/update/counter/PollCount/\d+\z`, responder)
	httpmock.ActivateNonDefault(c.Client.GetClient())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go c.Run(&wg)
	time.Sleep(2 * time.Second)
	c.Stop()
	wg.Wait()

	pollCount, err := storage.GetCounter(ctx, "PollCount")
	require.NoError(t, err)
	assert.Equal(t, int64(0), pollCount)
}
