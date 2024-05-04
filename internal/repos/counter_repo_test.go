package repos

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func BenchmarkCounterRepo_Update(b *testing.B) {
	c := NewCounterRepo()
	ctx := context.Background()
	v := rand.Int63()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Update(ctx, "Name", v)
	}
}

func BenchmarkCounterRepo_Set(b *testing.B) {
	c := NewCounterRepo()
	ctx := context.Background()
	v := rand.Int63()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set(ctx, "Name", v)
	}
}

func BenchmarkCounterRepo_Get(b *testing.B) {
	c := NewCounterRepo()
	ctx := context.Background()
	v := rand.Int63()
	c.Update(ctx, "Name", v)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Get(ctx, "Name")
	}
}

func BenchmarkCounterRepo_GetAll(b *testing.B) {
	c := NewCounterRepo()
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		c.Update(ctx, RandStringRunes(10), rand.Int63())
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.GetAll(ctx)
	}
}

func TestCounterRepo_Update(t *testing.T) {
	tests := []struct {
		errValue        error
		existingMetrics map[string]int64
		name            string
		metricName      string
		metricValue     int64
		wantValue       int64
		wantErr         bool
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
		errValue        error
		existingMetrics map[string]int64
		name            string
		metricName      string
		metricValue     int64
		wantValue       int64
		wantErr         bool
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
		errValue        error
		existingMetrics map[string]int64
		name            string
		metricName      string
		wantValue       int64
		wantErr         bool
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
		metrics map[string]int64
		want    map[string]int64
		name    string
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

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
