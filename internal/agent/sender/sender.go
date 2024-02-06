package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"net/http"
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
	s.sendGauges()
	s.sendCounters()
	ctx := context.Background()
	_, err := s.Storage.SetCounter(ctx, "PollCount", 0)
	if err != nil {
		logger.Log.Error(
			"Failed to set metric",
			zap.String("name", "PollCount"),
			zap.Int64("value", 0),
			zap.Error(err))
	}
}

func (s *Sender) sendGauges() {
	ctx := context.Background()
	m, err := s.Storage.GetAllGauge(ctx)
	if err != nil {
		logger.Log.Error("Failed to get gauge metrics", zap.Error(err))
	}
	url := fmt.Sprintf("http://%s/update/", s.ServerAddr)
	for key, value := range m {
		fields := []zap.Field{
			zap.String("name", key),
			zap.String("type", models.Gauge),
			zap.Float64("value", value),
		}
		metric := &models.Metrics{
			ID:    key,
			MType: models.Gauge,
			Value: &value,
		}
		var b bytes.Buffer
		data, err := json.Marshal(metric)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Log.Error("Failed to marshal data", fields...)
			continue
		}
		cw := gzip.NewWriter(&b)
		_, err = cw.Write(data)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Log.Error("Failed to compress data", fields...)
			continue
		}
		cw.Close()
		resp, err := s.Client.R().
			SetHeader("Content-Encoding", "gzip").
			SetBody(&b).
			Post(url)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Log.Error("Failed to send metric", fields...)
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			fields = append(fields, zap.Int("code", resp.StatusCode()))
			logger.Log.Error("Failed to send metric", fields...)
			continue
		}
		logger.Log.Info("Metric sent successfully", fields...)
	}
}

func (s *Sender) sendCounters() {
	ctx := context.Background()
	m, err := s.Storage.GetAllCounter(ctx)
	if err != nil {
		logger.Log.Error("Failed to get counter metrics", zap.Error(err))
	}
	url := fmt.Sprintf("http://%s/update/", s.ServerAddr)
	for key, value := range m {
		fields := []zap.Field{
			zap.String("name", key),
			zap.String("type", models.Counter),
			zap.Int64("value", value),
		}
		metric := &models.Metrics{
			ID:    key,
			MType: models.Counter,
			Delta: &value,
		}
		var b bytes.Buffer
		data, err := json.Marshal(metric)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Log.Error("Failed to marshal data", fields...)
			continue
		}
		cw := gzip.NewWriter(&b)
		_, err = cw.Write(data)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Log.Error("Failed to compress data", fields...)
			continue
		}
		cw.Close()
		resp, err := s.Client.R().
			SetHeader("Content-Encoding", "gzip").
			SetBody(&b).
			Post(url)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Log.Error("Failed to send metric", fields...)
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			fields = append(fields, zap.Int("code", resp.StatusCode()))
			logger.Log.Error("Failed to send metric", fields...)
			continue
		}
		logger.Log.Info("Metric sent successfully", fields...)
	}
}
