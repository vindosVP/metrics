package repos

import (
	"github.com/stretchr/testify/assert"
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
			val, err := repo.Update(tt.metricName, tt.metricValue)
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
			val, err := repo.Get(tt.metricName)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
			} else {
				assert.Equal(t, tt.wantValue, val)
			}
		})
	}
}
