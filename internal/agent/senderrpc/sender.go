package senderrpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/models"
	pb "github.com/vindosVP/metrics/internal/proto"
	"github.com/vindosVP/metrics/pkg/logger"
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
	Client         pb.MetricsClient
	ReportInterval time.Duration
	RateLimit      int
}

type job struct {
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
func New(cfg *config.AgentConfig, s MetricsStorage) *Sender {
	conn, err := grpc.NewClient(cfg.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal("Failed to connect GRPC")
	}
	return &Sender{
		Done:           make(chan struct{}),
		ReportInterval: cfg.ReportInterval,
		Storage:        s,
		Client:         pb.NewMetricsClient(conn),
		RateLimit:      cfg.RateLimit,
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
		err := s.send(j.metrics)
		results <- result{err, j.id}
	}
	wg.Done()
}

func (s *Sender) send(chunk []*models.Metrics) error {
	ctx := context.Background()
	_, err := retry.DoWithData(func() (*pb.UpdateBatchResponse, error) {
		metrics := make([]*pb.Metric, 0, len(chunk))
		for _, v := range chunk {
			if v.MType == models.Gauge {
				val := v.Value
				metrics = append(metrics, &pb.Metric{
					Type:  pb.MType_GAUGE,
					Id:    v.ID,
					Value: *val,
				})
			} else {
				val := v.Delta
				metrics = append(metrics, &pb.Metric{
					Type:  pb.MType_COUNTER,
					Id:    v.ID,
					Delta: *val,
				})
			}
		}
		return s.Client.UpdateBatch(ctx, &pb.UpdateBatchRequest{Metrics: metrics})
	}, retryOpts()...)

	if err != nil {
		return fmt.Errorf("failed to send metrics: %v", err)
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
		retry.OnRetry(func(n uint, _ error) {
			logger.Log.Info(fmt.Sprintf("Failed to connect to server, retrying in %s", retryDelays[n]))
		}),
		retry.Attempts(4),
		retry.LastErrorOnly(true),
	}
}
