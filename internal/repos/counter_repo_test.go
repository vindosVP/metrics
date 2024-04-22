package repos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounterRepo_Update(t *testing.T) {
	tests := []struct {
		name            string
		existingMetrics map[string]int64
		metricName      string
		metricValue     int64
		wantValue       int64
		wantErr         bool
		errValue        error
	}{
		{
			name:            "empty metrics",
			existingMetrics: make(map[string]int64),
			metricName:      "PollCount",
			metricValue:     123,
			wantValue:       123,
			wantErr:         false},
		{
			name: "existing metric",
			existingMetrics: map[string]int64{
				"PollCount": 1,
			},
			metricName:  "PollCount",
			metricValue: 123,
			wantValue:   124,
			wantErr:     false},
		{
			name: "existing negative metric",
			existingMetrics: map[string]int64{
				"PollCount": -5,
			},
			metricName:  "PollCount",
			metricValue: 123,
			wantValue:   118,
			wantErr:     false},
		{
			name:            "zero value",
			existingMetrics: make(map[string]int64),
			metricName:      "PollCount",
			metricValue:     0,
			wantValue:       0,
			wantErr:         false},
		{
			name:            "negative value",
			existingMetrics: make(map[string]int64),
			metricName:      "PollCount",
			metricValue:     -22,
			wantValue:       -22,
			wantErr:         false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CounterRepo{metrics: tt.existingMetrics}
			val, err := repo.Update(context.Background(), tt.metricName, tt.metricValue)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, val)
				assert.Equal(t, repo.metrics[tt.metricName], tt.wantValue)
			}
		})
	}

}

func TestCounterRepo_Set(t *testing.T) {
	tests := []struct {
		name            string
		existingMetrics map[string]int64
		metricName      string
		metricValue     int64
		wantValue       int64
		wantErr         bool
		errValue        error
	}{
		{
			name:            "empty metrics",
			existingMetrics: make(map[string]int64),
			metricName:      "PollCount",
			metricValue:     123,
			wantValue:       123,
			wantErr:         false},
		{
			name: "existing metric",
			existingMetrics: map[string]int64{
				"PollCount": 1,
			},
			metricName:  "PollCount",
			metricValue: 123,
			wantValue:   123,
			wantErr:     false},
		{
			name: "existing negative metric",
			existingMetrics: map[string]int64{
				"PollCount": -5,
			},
			metricName:  "PollCount",
			metricValue: 123,
			wantValue:   123,
			wantErr:     false},
		{
			name:            "zero value",
			existingMetrics: make(map[string]int64),
			metricName:      "PollCount",
			metricValue:     0,
			wantValue:       0,
			wantErr:         false},
		{
			name:            "negative value",
			existingMetrics: make(map[string]int64),
			metricName:      "PollCount",
			metricValue:     -22,
			wantValue:       -22,
			wantErr:         false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CounterRepo{metrics: tt.existingMetrics}
			val, err := repo.Set(context.Background(), tt.metricName, tt.metricValue)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, val)
				assert.Equal(t, repo.metrics[tt.metricName], tt.wantValue)
			}
		})
	}

}

func TestCounterRepo_Get(t *testing.T) {
	tests := []struct {
		name            string
		existingMetrics map[string]int64
		metricName      string
		wantValue       int64
		wantErr         bool
		errValue        error
	}{
		{
			name: "metric registered",
			existingMetrics: map[string]int64{
				"PollCount": 1,
			},
			metricName: "PollCount",
			wantValue:  1,
			wantErr:    false},
		{
			name:            "metric not registered",
			existingMetrics: make(map[string]int64),
			metricName:      "PollCount",
			wantErr:         true,
			errValue:        ErrMetricNotRegistered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CounterRepo{metrics: tt.existingMetrics}
			val, err := repo.Get(context.Background(), tt.metricName)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
			}
		})
	}
}

func TestCounterRepo_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		metrics map[string]int64
		want    map[string]int64
	}{
		{
			name:    "empty",
			metrics: make(map[string]int64),
			want:    make(map[string]int64),
		},
		{
			name: "filled",
			metrics: map[string]int64{
				"PollCount": 1,
				"Test":      2,
			},
			want: map[string]int64{
				"PollCount": 1,
				"Test":      2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CounterRepo{
				metrics: tt.metrics,
			}
			got, err := c.GetAll(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
