package repos

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGaugeRepo_Update(t *testing.T) {
	tests := []struct {
		name            string
		existingMetrics map[string]float64
		metricName      string
		metricValue     float64
		wantValue       float64
		wantErr         bool
		errValue        error
	}{
		{
			name:            "empty metrics",
			existingMetrics: make(map[string]float64),
			metricName:      "Alloc",
			metricValue:     1994.43,
			wantValue:       1994.43,
			wantErr:         false},
		{
			name: "existing metric",
			existingMetrics: map[string]float64{
				"Alloc": 1.74832,
			},
			metricName:  "Alloc",
			metricValue: 1994.43,
			wantValue:   1994.43,
			wantErr:     false},
		{
			name:            "zero value",
			existingMetrics: make(map[string]float64),
			metricName:      "Alloc",
			metricValue:     0,
			wantValue:       0,
			wantErr:         false},
		{
			name:            "negative value",
			existingMetrics: make(map[string]float64),
			metricName:      "Alloc",
			metricValue:     -9902.33,
			wantValue:       -9902.33,
			wantErr:         false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &GaugeRepo{metrics: tt.existingMetrics}
			val, err := repo.Update(context.Background(), tt.metricName, tt.metricValue)
			assert.Equal(t, tt.wantValue, val)
			assert.Equal(t, repo.metrics[tt.metricName], tt.wantValue)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestGaugeRepo_Get(t *testing.T) {
	tests := []struct {
		name            string
		existingMetrics map[string]float64
		metricName      string
		wantValue       float64
		wantErr         bool
		errValue        error
	}{
		{
			name: "metric registered",
			existingMetrics: map[string]float64{
				"Alloc": 1,
			},
			metricName: "Alloc",
			wantValue:  1,
			wantErr:    false},
		{
			name:            "metric not registered",
			existingMetrics: make(map[string]float64),
			metricName:      "Alloc",
			wantErr:         true,
			errValue:        ErrMetricNotRegistered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &GaugeRepo{metrics: tt.existingMetrics}
			val, err := repo.Get(context.Background(), tt.metricName)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
			}
		})
	}
}

func TestGaugeRepo_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		metrics map[string]float64
		want    map[string]float64
	}{
		{
			name:    "empty",
			metrics: make(map[string]float64),
			want:    make(map[string]float64),
		},
		{
			name: "filled",
			metrics: map[string]float64{
				"Test":  2.2222222,
				"Test2": 0.000000000000000003,
			},
			want: map[string]float64{
				"Test":  2.2222222,
				"Test2": 0.000000000000000003,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GaugeRepo{
				metrics: tt.metrics,
			}
			got, err := g.GetAll(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
