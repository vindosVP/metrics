package memstorage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vindosVP/metrics/internal/storage/memstorage/mocks"
)

func TestStorage_UpdateCounter(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		mockErr     error
		errValue    error
		name        string
		metricName  string
		mockValue   int64
		metricValue int64
		wantValue   int64
		wantErr     bool
	}{
		{
			name:        "ok",
			mockValue:   12,
			mockErr:     nil,
			metricName:  "PollCount",
			metricValue: 12,
			wantValue:   12,
			wantErr:     false,
			errValue:    nil},
		{
			name:        "error",
			mockValue:   12,
			mockErr:     unexpectedError,
			metricName:  "PollCount",
			metricValue: 12,
			wantErr:     true,
			errValue:    unexpectedError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCounter := mocks.NewCounter(t)
			mockGauge := mocks.NewGauge(t)
			storage := New(mockGauge, mockCounter)
			ctx := context.Background()
			mockCounter.On("Update", mock.Anything, tt.metricName, tt.metricValue).Return(tt.mockValue, tt.mockErr)
			val, err := storage.UpdateCounter(ctx, tt.metricName, tt.metricValue)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
				assert.NoError(t, err, tt.wantErr)
			}
		})
	}

}

func TestStorage_SetCounter(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		mockErr     error
		errValue    error
		name        string
		metricName  string
		mockValue   int64
		metricValue int64
		wantValue   int64
		wantErr     bool
	}{
		{
			name:        "ok",
			mockValue:   12,
			mockErr:     nil,
			metricName:  "PollCount",
			metricValue: 12,
			wantValue:   12,
			wantErr:     false,
			errValue:    nil},
		{
			name:        "error",
			mockValue:   12,
			mockErr:     unexpectedError,
			metricName:  "PollCount",
			metricValue: 12,
			wantErr:     true,
			errValue:    unexpectedError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCounter := mocks.NewCounter(t)
			mockGauge := mocks.NewGauge(t)
			storage := New(mockGauge, mockCounter)
			ctx := context.Background()
			mockCounter.On("Set", mock.Anything, tt.metricName, tt.metricValue).Return(tt.mockValue, tt.mockErr)
			val, err := storage.SetCounter(ctx, tt.metricName, tt.metricValue)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
				assert.NoError(t, err, tt.wantErr)
			}
		})
	}

}

func TestStorage_UpdateGauge(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		mockErr     error
		errValue    error
		name        string
		metricName  string
		mockValue   float64
		metricValue float64
		wantValue   float64
		wantErr     bool
	}{
		{
			name:        "ok",
			mockValue:   0.000003,
			mockErr:     nil,
			metricName:  "Alloc",
			metricValue: 0.000003,
			wantValue:   0.000003,
			wantErr:     false,
			errValue:    nil},
		{
			name:        "error",
			mockValue:   0.000003,
			mockErr:     unexpectedError,
			metricName:  "Alloc",
			metricValue: 0.000003,
			wantErr:     true,
			errValue:    unexpectedError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCounter := mocks.NewCounter(t)
			mockGauge := mocks.NewGauge(t)
			storage := New(mockGauge, mockCounter)
			ctx := context.Background()
			mockGauge.On("Update", mock.Anything, tt.metricName, tt.metricValue).Return(tt.mockValue, tt.mockErr)
			val, err := storage.UpdateGauge(ctx, tt.metricName, tt.metricValue)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
				assert.NoError(t, err, tt.wantErr)
			}
		})
	}

}

func TestStorage_GetGauge(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		mockErr    error
		errValue   error
		name       string
		metricName string
		mockValue  float64
		wantValue  float64
		wantErr    bool
	}{
		{
			name:       "ok",
			mockValue:  0.000003,
			mockErr:    nil,
			metricName: "Alloc",
			wantValue:  0.000003,
			wantErr:    false,
			errValue:   nil},
		{
			name:       "error",
			mockValue:  0.000003,
			mockErr:    unexpectedError,
			metricName: "Alloc",
			wantValue:  0,
			wantErr:    true,
			errValue:   unexpectedError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCounter := mocks.NewCounter(t)
			mockGauge := mocks.NewGauge(t)
			storage := New(mockGauge, mockCounter)
			ctx := context.Background()
			mockGauge.On("Get", mock.Anything, tt.metricName).Return(tt.mockValue, tt.mockErr)
			val, err := storage.GetGauge(ctx, tt.metricName)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
				assert.NoError(t, err, tt.wantErr)
			}
		})
	}

}

func TestStorage_GetCounter(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		mockErr    error
		errValue   error
		name       string
		metricName string
		mockValue  int64
		wantValue  int64
		wantErr    bool
	}{
		{
			name:       "ok",
			mockValue:  15,
			mockErr:    nil,
			metricName: "PollCount",
			wantValue:  15,
			wantErr:    false,
			errValue:   nil},
		{
			name:       "error",
			mockValue:  15,
			mockErr:    unexpectedError,
			metricName: "PollCount",
			wantValue:  0,
			wantErr:    true,
			errValue:   unexpectedError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCounter := mocks.NewCounter(t)
			mockGauge := mocks.NewGauge(t)
			storage := New(mockGauge, mockCounter)
			ctx := context.Background()
			mockCounter.On("Get", mock.Anything, tt.metricName).Return(tt.mockValue, tt.mockErr)
			val, err := storage.GetCounter(ctx, tt.metricName)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
				assert.NoError(t, err, tt.wantErr)
			}
		})
	}

}

func TestStorage_GetAllGauge(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		errValue    error
		mockMetrics map[string]float64
		wantMetrics map[string]float64
		name        string
		wantErr     bool
	}{
		{
			name:        "empty",
			mockMetrics: make(map[string]float64),
			wantMetrics: make(map[string]float64),
			wantErr:     false,
			errValue:    nil,
		},
		{
			name: "filled",
			mockMetrics: map[string]float64{
				"Alloc": 12.2,
				"Test":  33,
			},
			wantMetrics: map[string]float64{
				"Alloc": 12.2,
				"Test":  33,
			},
			wantErr:  false,
			errValue: nil,
		},
		{
			name:        "unexpected error",
			mockMetrics: nil,
			wantMetrics: nil,
			wantErr:     true,
			errValue:    unexpectedError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCounter := mocks.NewCounter(t)
			mockGauge := mocks.NewGauge(t)
			storage := New(mockGauge, mockCounter)
			ctx := context.Background()
			mockGauge.On("GetAll", mock.Anything).Return(tt.mockMetrics, tt.errValue)

			got, err := storage.GetAllGauge(ctx)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMetrics, got)
			}
		})
	}
}

func TestStorage_GetAllCounter(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		errValue    error
		mockMetrics map[string]int64
		wantMetrics map[string]int64
		name        string
		wantErr     bool
	}{
		{
			name:        "empty",
			mockMetrics: make(map[string]int64),
			wantMetrics: make(map[string]int64),
			wantErr:     false,
			errValue:    nil,
		},
		{
			name: "filled",
			mockMetrics: map[string]int64{
				"PollCount": 35,
				"Test":      33,
			},
			wantMetrics: map[string]int64{
				"PollCount": 35,
				"Test":      33,
			},
			wantErr:  false,
			errValue: nil,
		},
		{
			name:        "unexpected error",
			mockMetrics: nil,
			wantMetrics: nil,
			wantErr:     true,
			errValue:    unexpectedError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCounter := mocks.NewCounter(t)
			mockGauge := mocks.NewGauge(t)
			storage := New(mockGauge, mockCounter)
			mockCounter.On("GetAll", mock.Anything).Return(tt.mockMetrics, tt.errValue)

			got, err := storage.GetAllCounter(context.Background())
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMetrics, got)
			}
		})
	}
}
