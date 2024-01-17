package sender

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vindosVP/metrics/cmd/agent/config"
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
		log.Printf("Failed to set PollCount to 0: %v", err)
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
			log.Printf("Failed to send %s:%v", key, err)
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			log.Printf("Failed to send %s: resp code: %d, data %s", key, resp.StatusCode(), string(resp.Body()))
			continue
		}
		log.Printf("Metric %s sent sucessfully", key)
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
			log.Printf("Failed to send %s:%v", key, err)
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			log.Printf("Failed to send %s: resp code: %d, data %s", key, resp.StatusCode(), string(resp.Body()))
			continue
		}
		log.Printf("Metric %s sent sucessfully", key)
	}
}
