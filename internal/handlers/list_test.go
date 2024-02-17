package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vindosVP/metrics/internal/handlers/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(t *testing.T) {
	type mockGauge struct {
		needed bool
		fields map[string]float64
		err    error
	}
	type mockCounter struct {
		needed bool
		fields map[string]int64
		err    error
	}
	type want struct {
		code        int
		contentType string
	}
	unexpectedError := errors.New("unexpected error")

	tests := []struct {
		name        string
		mockGauge   mockGauge
		mockCounter mockCounter
		method      string
		want        want
	}{
		{
			name: "gauge error",
			mockGauge: mockGauge{
				needed: true,
				fields: nil,
				err:    unexpectedError,
			},
			mockCounter: mockCounter{
				needed: true,
				fields: make(map[string]int64),
				err:    nil,
			},
			method: http.MethodGet,
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "counter error",
			mockGauge: mockGauge{
				needed: false,
				fields: nil,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: true,
				fields: nil,
				err:    unexpectedError,
			},
			method: http.MethodGet,
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong method",
			mockGauge: mockGauge{
				needed: false,
				fields: nil,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: false,
				fields: nil,
				err:    unexpectedError,
			},
			method: http.MethodPost,
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name: "ok",
			mockGauge: mockGauge{
				needed: true,
				fields: make(map[string]float64),
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: true,
				fields: make(map[string]int64),
				err:    nil,
			},
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewMetricsStorage(t)
			if tt.mockGauge.needed {
				mockStorage.On("GetAllGauge", mock.Anything).Return(tt.mockGauge.fields, tt.mockGauge.err)
			}
			if tt.mockCounter.needed {
				mockStorage.On("GetAllCounter", mock.Anything).Return(tt.mockCounter.fields, tt.mockCounter.err)
			}

			r := chi.NewRouter()
			r.Get("/", List(mockStorage))

			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			defer res.Body.Close()
		})
	}
}

func Test_counterMetricLines(t *testing.T) {

	tests := []struct {
		name    string
		metrics map[string]int64
		want    []string
	}{
		{
			name: "filled",
			metrics: map[string]int64{
				"PollCount": 12,
				"Test":      1,
			},
			want: []string{"<tr><td>PollCount</td><td>12</td></tr>", "<tr><td>Test</td><td>1</td></tr>"},
		},
		{
			name:    "empty",
			metrics: make(map[string]int64),
			want:    make([]string, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := counterMetricLines(tt.metrics)
			assert.ElementsMatch(t, lines, tt.want)
		})
	}

}

func Test_gaugeMetricLines(t *testing.T) {

	tests := []struct {
		name    string
		metrics map[string]float64
		want    []string
	}{
		{
			name: "filled",
			metrics: map[string]float64{
				"Alloc": 323423452.555,
				"Test":  1,
			},
			want: []string{"<tr><td>Alloc</td><td>323423452.56</td></tr>", "<tr><td>Test</td><td>1.00</td></tr>"},
		},
		{
			name:    "empty",
			metrics: make(map[string]float64),
			want:    make([]string, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := gaugeMetricLines(tt.metrics)
			assert.ElementsMatch(t, lines, tt.want)
		})
	}

}
