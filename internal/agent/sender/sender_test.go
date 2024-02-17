package sender

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"net/http"
	"testing"
	"time"
)

func TestSender(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test because testing.Short is enabled")
	}

	done := make(chan struct{})
	senderShutdown := make(chan struct{})

	cfg := config.NewAgentConfig()
	cRepo := repos.NewCounterRepo()
	ctx := context.Background()
	_, err := cRepo.Set(ctx, "PollCount", 100)
	require.NoError(t, err)
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)
	c := New(cfg, storage)
	c.Done = done
	c.ReportInterval = 1

	responder := httpmock.NewStringResponder(200, "")
	httpmock.RegisterResponder(http.MethodPost, `=~^(http|https)://.+/update/counter/PollCount/\d+\z`, responder)
	httpmock.ActivateNonDefault(c.Client.GetClient())

	go func() {
		defer close(senderShutdown)
		c.Run()
	}()
	time.Sleep(2 * time.Second)
	close(done)
	<-senderShutdown

	pollCount, err := storage.GetCounter(ctx, "PollCount")
	require.NoError(t, err)
	assert.Equal(t, int64(0), pollCount)
}
