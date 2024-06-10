package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/metrics/internal/models"
	pb "github.com/vindosVP/metrics/internal/proto"
	"github.com/vindosVP/metrics/internal/storage"
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

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	s MetricsStorage
}

func NewMetricsServer(s MetricsStorage) *MetricsServer {
	return &MetricsServer{s: s}
}

func (s *MetricsServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {

	var resp pb.GetResponse

	fields := []zap.Field{
		zap.String("name", in.Id),
		zap.String("type", in.Type.String()),
	}

	switch in.Type {
	case pb.MType_COUNTER:
		val, cerr := s.s.GetCounter(ctx, in.Id)
		if cerr != nil {
			code := codes.Internal
			if errors.Is(cerr, storage.ErrMetricNotRegistered) {
				code = codes.NotFound
			} else {
				fields = append(fields, zap.Error(cerr))
				logger.Log.Error("Failed to get metric value", fields...)
			}
			return &resp, status.Errorf(code, "failed to get counter: %s", cerr)
		}

		resp.Metric.Id = in.Id
		resp.Metric.Type = pb.MType_COUNTER
		resp.Metric.Delta = val
	case pb.MType_GAUGE:
		val, gerr := s.s.GetGauge(ctx, in.Id)
		if gerr != nil {
			code := codes.Internal
			if errors.Is(gerr, storage.ErrMetricNotRegistered) {
				code = codes.NotFound
			} else {
				fields = append(fields, zap.Error(gerr))
				logger.Log.Error("Failed to get metric value", fields...)
			}
			return &resp, status.Errorf(code, "failed to get gauge: %s", gerr)
		}

		resp.Metric.Id = in.Id
		resp.Metric.Type = pb.MType_GAUGE
		resp.Metric.Value = val
	}

	return &resp, nil
}

func (s *MetricsServer) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	var resp pb.UpdateResponse
	resp.Metric = &pb.Metric{}

	fields := []zap.Field{
		zap.String("name", in.Metric.Id),
		zap.String("type", in.Metric.Type.String()),
	}

	switch in.Metric.Type {
	case pb.MType_COUNTER:

		delta := in.Metric.Delta
		fields = append(fields, zap.Int64("delta", delta))
		val, cerr := s.s.UpdateCounter(ctx, in.Metric.Id, delta)
		if cerr != nil {
			fields = append(fields, zap.Error(cerr))
			logger.Log.Error("Failed to update metric value", fields...)
			return nil, status.Errorf(codes.Internal, "failed to update metric: %s", cerr)
		}

		resp.Metric.Id = in.Metric.Id
		resp.Metric.Type = pb.MType_COUNTER
		resp.Metric.Delta = val

		logger.Log.Info("Updated metric value", fields...)
	case pb.MType_GAUGE:
		value := in.Metric.Value
		fields = append(fields, zap.Float64("value", value))
		val, gerr := s.s.UpdateGauge(ctx, in.Metric.Id, value)
		if gerr != nil {
			fields = append(fields, zap.Error(gerr))
			logger.Log.Error("Failed to update metric value", fields...)
			return nil, status.Errorf(codes.Internal, "failed to update metric: %s", gerr)
		}

		resp.Metric.Id = in.Metric.Id
		resp.Metric.Type = pb.MType_GAUGE
		resp.Metric.Value = val

		logger.Log.Info("Updated metric value", fields...)
	}

	return &resp, nil
}

func (s *MetricsServer) UpdateBatch(ctx context.Context, in *pb.UpdateBatchRequest) (*pb.UpdateBatchResponse, error) {

	var resp pb.UpdateBatchResponse

	batch := make([]*models.Metrics, 0)
	for _, v := range in.Metrics {
		var mType string
		if v.Type == pb.MType_COUNTER {
			mType = models.Counter
		} else {
			mType = models.Gauge
		}
		metric := &models.Metrics{
			Delta: &v.Delta,
			Value: &v.Value,
			ID:    v.Id,
			MType: mType,
		}
		batch = append(batch, metric)
	}

	err := s.s.InsertBatch(ctx, batch)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert batch: %s", err)
	}

	return &resp, nil
}
