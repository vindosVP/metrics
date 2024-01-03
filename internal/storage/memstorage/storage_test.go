package memstorage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vindosVP/metrics/internal/repos/mocks"
	"testing"
)

func TestStorage_UpdateCounter(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	tests := []struct {
		name        string
		mockValue   int64
		mockErr     error
		metricName  string
		metricValue int64
		wantValue   int64
		wantErr     bool
		errValue    error
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
			mockCounter.On("Update", tt.metricName, tt.metricValue).Return(tt.mockValue, tt.mockErr)
			val, err := storage.UpdateCounter(tt.metricName, tt.metricValue)

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
		name        string
		mockValue   int64
		mockErr     error
		metricName  string
		metricValue int64
		wantValue   int64
		wantErr     bool
		errValue    error
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
			mockCounter.On("Set", tt.metricName, tt.metricValue).Return(tt.mockValue, tt.mockErr)
			val, err := storage.SetCounter(tt.metricName, tt.metricValue)

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
		name        string
		mockValue   float64
		mockErr     error
		metricName  string
		metricValue float64
		wantValue   float64
		wantErr     bool
		errValue    error
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
			mockGauge.On("Update", tt.metricName, tt.metricValue).Return(tt.mockValue, tt.mockErr)
			val, err := storage.UpdateGauge(tt.metricName, tt.metricValue)

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
		name       string
		mockValue  float64
		mockErr    error
		metricName string
		wantValue  float64
		wantErr    bool
		errValue   error
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
			mockGauge.On("Get", tt.metricName).Return(tt.mockValue, tt.mockErr)
			val, err := storage.GetGauge(tt.metricName)

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
		name       string
		mockValue  int64
		mockErr    error
		metricName string
		wantValue  int64
		wantErr    bool
		errValue   error
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
			mockCounter.On("Get", tt.metricName).Return(tt.mockValue, tt.mockErr)
			val, err := storage.GetCounter(tt.metricName)

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
		name        string
		mockMetrics map[string]float64
		wantMetrics map[string]float64
		wantErr     bool
		errValue    error
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
			mockGauge.On("GetAll").Return(tt.mockMetrics, tt.errValue)

			got, err := storage.GetAllGauge()
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
		name        string
		mockMetrics map[string]int64
		wantMetrics map[string]int64
		wantErr     bool
		errValue    error
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
			mockCounter.On("GetAll").Return(tt.mockMetrics, tt.errValue)

			got, err := storage.GetAllCounter()
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMetrics, got)
			}
		})
	}
}
