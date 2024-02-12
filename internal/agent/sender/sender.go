package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"
	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"syscall"
	"time"
)

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

type Sender struct {
	ReportInterval time.Duration
	ServerAddr     string
	Done           <-chan struct{}
	Storage        MetricsStorage
	Client         *resty.Client
}

func New(cfg *config.AgentConfig, s MetricsStorage) *Sender {
	return &Sender{
		ReportInterval: cfg.ReportInterval,
		ServerAddr:     cfg.ServerAddr,
		Storage:        s,
		Client:         resty.New(),
	}
}

var retryDelays = map[uint]time.Duration{
	0: 1 * time.Second,
	1: 3 * time.Second,
	2: 5 * time.Second,
}

func (s *Sender) Run() {
	tick := time.NewTicker(s.ReportInterval * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-s.Done:
			return
		case <-tick.C:
			s.SendMetrics()
		}
	}
}

func (s *Sender) SendMetrics() {
	ctx := context.Background()
	g, err := s.Storage.GetAllGauge(ctx)
	if err != nil {
		logger.Log.Error("Failed to get gauge metrics", zap.Error(err))
	}
	c, err := s.Storage.GetAllCounter(ctx)
	if err != nil {
		logger.Log.Error("Failed to get counter metrics", zap.Error(err))
	}

	s.send(c, g)
	_, err = s.Storage.SetCounter(ctx, "PollCount", 0)
	if err != nil {
		logger.Log.Error(
			"Failed to set metric",
			zap.String("name", "PollCount"),
			zap.Int64("value", 0),
			zap.Error(err))
	}
}

func (s *Sender) send(c map[string]int64, g map[string]float64) {
	if len(c)+len(g) == 0 {
		return
	}

	batch := makeButch(c, g)

	var b bytes.Buffer
	data, err := json.Marshal(batch)
	if err != nil {
		logger.Log.Error("Failed to marshal data", zap.Error(err))
		return
	}

	cw := gzip.NewWriter(&b)
	_, err = cw.Write(data)
	if err != nil {
		logger.Log.Error("Failed to compress data", zap.Error(err))
		return
	}
	cw.Close()

	url := fmt.Sprintf("http://%s/updates/", s.ServerAddr)
	resp, err := retry.DoWithData(func() (*resty.Response, error) {
		return s.Client.R().
			SetHeader("Content-Encoding", "gzip").
			SetBody(&b).
			Post(url)
	}, retryOpts()...)

	if err != nil {
		logger.Log.Error("Failed to send metrics", zap.Error(err))
		return
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Error("Failed to send metrics", zap.Int("code", resp.StatusCode()))
		return
	}
	logger.Log.Info("Metric sent successfully")
}

func makeButch(c map[string]int64, g map[string]float64) []*models.Metrics {
	batch := make([]*models.Metrics, len(c)+len(g))

	i := 0
	for k, v := range g {
		metric := &models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &v,
		}
		batch[i] = metric
		i++
	}

	for k, v := range c {
		metric := &models.Metrics{
			ID:    k,
			MType: models.Counter,
			Delta: &v,
		}
		batch[i] = metric
		i++
	}

	return batch
}

func retryOpts() []retry.Option {
	return []retry.Option{
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, syscall.ECONNREFUSED)
		}),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			delay := retryDelays[n]
			return delay
		}),
		retry.OnRetry(func(n uint, err error) {
			logger.Log.Info(fmt.Sprintf("Failed to connect to server, retrying in %s", retryDelays[n]))
		}),
		retry.Attempts(4),
	}
}
