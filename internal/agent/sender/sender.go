// Package sender sends collected metrics to the server every n seconds
package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/pkg/encryption"
	"github.com/vindosVP/metrics/pkg/logger"
	"github.com/vindosVP/metrics/pkg/utils"
)

// MetricsStorage consists of methods to write and get data from storage.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=MetricsStorage
type MetricsStorage interface {

	// SetCounter sets counter metric value.
	// new value replaces the old one.
	SetCounter(ctx context.Context, name string, v int64) (int64, error)

	// GetAllGauge returns all collected gauge metrics.
	GetAllGauge(ctx context.Context) (map[string]float64, error)

	// GetAllCounter returns all collected counter metrics.
	GetAllCounter(ctx context.Context) (map[string]int64, error)
}

const chunkSize = 3

// Sender consists data to send metrics
type Sender struct {
	Storage        MetricsStorage
	Done           chan struct{}
	Client         *resty.Client
	ServerAddr     string
	Key            string
	ReportInterval time.Duration
	RateLimit      int
	UseHash        bool
	CryptoKey      *rsa.PublicKey
	IP             net.IP
}

type job struct {
	url     string
	metrics []*models.Metrics
	id      int
}

type result struct {
	err error
	id  int
}

var retryDelays = map[uint]time.Duration{
	0: 1 * time.Second,
	1: 3 * time.Second,
	2: 5 * time.Second,
}

// New creates the Sender
func New(cfg *config.AgentConfig, s MetricsStorage, cryptoKey *rsa.PublicKey, IP net.IP) *Sender {
	return &Sender{
		Done:           make(chan struct{}),
		ReportInterval: cfg.ReportInterval,
		ServerAddr:     cfg.ServerAddr,
		Storage:        s,
		Client:         resty.New(),
		UseHash:        cfg.Key != "",
		Key:            cfg.Key,
		RateLimit:      cfg.RateLimit,
		CryptoKey:      cryptoKey,
		IP:             IP,
	}
}

func (s *Sender) Stop() {
	close(s.Done)
}

// Run starts the Sender to send metrics
func (s *Sender) Run(wg *sync.WaitGroup) {
	tick := time.NewTicker(s.ReportInterval * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-s.Done:
			wg.Done()
			return
		case <-tick.C:
			s.sendMetrics()
		}
	}
}

func (s *Sender) sendMetrics() {
	ctx := context.Background()
	g, err := s.Storage.GetAllGauge(ctx)
	if err != nil {
		logger.Log.Error("Failed to get gauge metrics", zap.Error(err))
	}
	c, err := s.Storage.GetAllCounter(ctx)
	if err != nil {
		logger.Log.Error("Failed to get counter metrics", zap.Error(err))
	}
	batch := makeButch(c, g)
	jobs := s.generateJobs(batch)
	results := make(chan result)
	go listenResults(results)
	s.startWorkers(jobs, results, s.RateLimit)
	_, err = s.Storage.SetCounter(ctx, "PollCount", 0)
	if err != nil {
		logger.Log.Error(
			"Failed to set metric",
			zap.String("name", "PollCount"),
			zap.Int64("value", 0),
			zap.Error(err))
	}
}

func (s *Sender) generateJobs(metrics []*models.Metrics) chan job {
	jobs := make(chan job)
	go func() {
		size := chunkSize
		url := fmt.Sprintf("http://%s/updates/", s.ServerAddr)
		id := 1
		for {
			if len(metrics) == 0 {
				break
			}
			if len(metrics) < size {
				size = len(metrics)
			}
			jobs <- job{
				id:      id,
				url:     url,
				metrics: metrics[0:size],
			}
			metrics = metrics[size:]
			id++
		}
		defer close(jobs)
	}()
	return jobs
}

func listenResults(results <-chan result) {
	for res := range results {
		if res.err != nil {
			logger.Log.Error("worker failed", zap.Error(res.err), zap.Int("id", res.id))
		} else {
			logger.Log.Info("worker finished", zap.Int("id", res.id))
		}
	}
}

func (s *Sender) startWorkers(jobs <-chan job, results chan<- result, workers int) {
	wg := sync.WaitGroup{}
	logger.Log.Info(fmt.Sprintf("Starting %d workers", workers))
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go s.worker(jobs, results, &wg)
	}
	wg.Wait()
	close(results)
}

func (s *Sender) worker(jobs <-chan job, results chan<- result, wg *sync.WaitGroup) {
	for j := range jobs {
		err := s.send(j.url, j.metrics, s.IP)
		results <- result{err, j.id}
	}
	wg.Done()
}

func (s *Sender) send(url string, chunk []*models.Metrics, ip net.IP) error {

	var b bytes.Buffer
	data, err := json.Marshal(chunk)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %v", err)
	}

	cw := gzip.NewWriter(&b)
	_, err = cw.Write(data)
	if err != nil {
		return fmt.Errorf("failed to gzip metrics: %v", err)
	}
	cw.Close()

	hash := ""
	if s.UseHash {
		hash, err = utils.Sha256Hash(b.Bytes(), s.Key)
		if err != nil {
			return fmt.Errorf("failed to hash metrics: %v", err)
		}
	}
	var body []byte
	if s.CryptoKey != nil {
		body, err = encryption.Encrypt(s.CryptoKey, b.Bytes())
		if err != nil {
			return fmt.Errorf("failed to encrypt metrics: %v", err)
		}
	} else {
		body = b.Bytes()
	}

	resp, err := retry.DoWithData(func() (*resty.Response, error) {
		req := s.Client.R().
			SetHeader("Content-Encoding", "gzip").
			SetHeader("X-Real-IP", ip.String()).
			SetBody(body)
		if s.UseHash {
			req.SetHeader("HashSHA256", hash)
		}
		return req.Post(url)
	}, retryOpts()...)

	if err != nil {
		return fmt.Errorf("failed to send metrics: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to send metrics: %v", resp.Status())
	}
	return nil
}

func makeButch(c map[string]int64, g map[string]float64) []*models.Metrics {
	batch := make([]*models.Metrics, len(c)+len(g))

	i := 0
	for k, v := range g {
		val := v
		metric := &models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &val,
		}
		batch[i] = metric
		i++
	}

	for k, v := range c {
		val := v
		metric := &models.Metrics{
			ID:    k,
			MType: models.Counter,
			Delta: &val,
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
		retry.LastErrorOnly(true),
	}
}
