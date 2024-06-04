package grpcserver

import (
	"context"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/vindosVP/metrics/internal/models"
	pb "github.com/vindosVP/metrics/internal/proto"
	"github.com/vindosVP/metrics/internal/service"
	"github.com/vindosVP/metrics/pkg/logger"
)

type MetricsStorage interface {
	UpdateGauge(ctx context.Context, name string, v float64) (float64, error)
	UpdateCounter(ctx context.Context, name string, v int64) (int64, error)
	SetCounter(ctx context.Context, name string, v int64) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
	GetAllGauge(ctx context.Context) (map[string]float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAllCounter(ctx context.Context) (map[string]int64, error)
	InsertBatch(ctx context.Context, batch []*models.Metrics) error
}

type GRPCServer struct {
	listen net.Listener
	s      *grpc.Server
}

func (g *GRPCServer) Run(wg *sync.WaitGroup) {
	err := g.s.Serve(g.listen)
	if err != nil {
		wg.Done()
		logger.Log.Error("failed to start GRPC server", zap.Error(err))
	}
}

func (g *GRPCServer) Stop(wg *sync.WaitGroup) {
	g.s.GracefulStop()
	wg.Done()
}

func New(st MetricsStorage, addr string) (*GRPCServer, error) {
	a := grpc.NewServer()
	pb.RegisterMetricsServer(a, service.NewMetricsServer(st))
	listen, err := net.Listen("tcp", addr)
	logger.Log.Info(fmt.Sprintf("GRPC server listening on %s", addr))
	if err != nil {
		return nil, fmt.Errorf("failed to create GRPC server: %w", err)
	}
	return &GRPCServer{
		listen: listen,
		s:      a,
	}, nil
}
