package sender

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {
	UpdateGauge(name string, v float64) (float64, error)
	UpdateCounter(name string, v int64) (int64, error)
	SetCounter(name string, v int64) (int64, error)
	GetGauge(name string) (float64, error)
	GetAllGauge() (map[string]float64, error)
	GetCounter(name string) (int64, error)
	GetAllCounter() (map[string]int64, error)
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
	_, err := s.Storage.SetCounter("PollCount", 0)
	if err != nil {
		logger.Log.Error(
			"Failed to set metric",
			zap.String("name", "PollCount"),
			zap.Int64("value", 0),
			zap.Error(err))
	}
}

func (s *Sender) sendGauges() {
	m, err := s.Storage.GetAllGauge()
	if err != nil {
		log.Print("Failed to get gauge metrics")
	}
	for key, value := range m {
		url := fmt.Sprintf("http://%s/update/gauge/%s/%f", s.ServerAddr, key, value)
		resp, err := s.Client.R().Post(url)
		if err != nil {
			logger.Log.Error(
				"Failed to send metric",
				zap.String("name", key),
				zap.Error(err))
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			logger.Log.Error(
				"Failed to send metric",
				zap.String("name", key),
				zap.Int("code", resp.StatusCode()),
				zap.String("data", string(resp.Body())))
			continue
		}
		logger.Log.Info(
			"Metric sent successfully",
			zap.String("name", key))
	}
}

func (s *Sender) sendCounters() {
	m, err := s.Storage.GetAllCounter()
	if err != nil {
		log.Print("Failed to get counter metrics")
	}
	for key, value := range m {
		url := fmt.Sprintf("http://%s/update/counter/%s/%d", s.ServerAddr, key, value)
		resp, err := s.Client.R().Post(url)
		if err != nil {
			logger.Log.Error(
				"Failed to send metric",
				zap.String("name", key),
				zap.Error(err))
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			logger.Log.Error(
				"Failed to send metric",
				zap.String("name", key),
				zap.Int("code", resp.StatusCode()),
				zap.String("data", string(resp.Body())))
			continue
		}
		logger.Log.Info(
			"Metric sent successfully",
			zap.String("name", key))
	}
}
