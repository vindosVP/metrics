package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/metrics/internal/handlers/mocks"
	"github.com/vindosVP/metrics/internal/storage"
)

func TestGetBody(t *testing.T) {
	type mockGauge struct {
		needed bool
		name   string
		value  float64
		err    error
	}
	type mockCounter struct {
		needed bool
		name   string
		value  int64
		err    error
	}
	type want struct {
		code        int
		body        string
		contentType string
	}
	unexpectedError := errors.New("unexpected error")

	tests := []struct {
		name        string
		mockGauge   mockGauge
		mockCounter mockCounter
		body        string
		method      string
		want        want
	}{
		{
			name: "wrong metric type",
			mockGauge: mockGauge{
				needed: false,
				name:   "",
				value:  1,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: false,
				name:   "",
				value:  1,
				err:    nil,
			},
			body:   "{\"id\":\"Alloc\",\"type\":\"WrongType\"}",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong method",
			mockGauge: mockGauge{
				needed: false,
				name:   "",
				value:  1,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: false,
				name:   "",
				value:  1,
				err:    nil,
			},
			body:   "{\"id\":\"Alloc\",\"type\":\"gauge\"}",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "counter unexpected error",
			mockGauge: mockGauge{
				needed: false,
				name:   "",
				value:  1,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: true,
				name:   "test",
				value:  0,
				err:    unexpectedError,
			},
			body:   "{\"id\":\"test\",\"type\":\"counter\"}",
			method: http.MethodPost,
			want: want{
				code:        http.StatusInternalServerError,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "gauge unexpected error",
			mockGauge: mockGauge{
				needed: true,
				name:   "test",
				value:  0,
				err:    unexpectedError,
			},
			mockCounter: mockCounter{
				needed: false,
				name:   "",
				value:  0,
				err:    nil,
			},
			body:   "{\"id\":\"test\",\"type\":\"gauge\"}",
			method: http.MethodPost,
			want: want{
				code:        http.StatusInternalServerError,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "gauge not registered",
			mockGauge: mockGauge{
				needed: true,
				name:   "test",
				value:  0,
				err:    storage.ErrMetricNotRegistered,
			},
			mockCounter: mockCounter{
				needed: false,
				name:   "",
				value:  0,
				err:    nil,
			},
			body:   "{\"id\":\"test\",\"type\":\"gauge\"}",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "counter not registered",
			mockGauge: mockGauge{
				needed: false,
				name:   "",
				value:  0,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: true,
				name:   "test",
				value:  0,
				err:    storage.ErrMetricNotRegistered,
			},
			body:   "{\"id\":\"test\",\"type\":\"counter\"}",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "gauge ok",
			mockGauge: mockGauge{
				needed: true,
				name:   "test",
				value:  122.44,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: false,
				name:   "",
				value:  0,
				err:    nil,
			},
			body:   "{\"id\":\"test\",\"type\":\"gauge\"}",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				body:        "{\"id\":\"test\",\"type\":\"gauge\",\"value\":122.44}",
				contentType: "application/json",
			},
		},
		{
			name: "counter ok",
			mockGauge: mockGauge{
				needed: false,
				name:   "",
				value:  0,
				err:    nil,
			},
			mockCounter: mockCounter{
				needed: true,
				name:   "test",
				value:  111,
				err:    nil,
			},
			body:   "{\"id\":\"test\",\"type\":\"counter\"}",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				body:        "{\"id\":\"test\",\"type\":\"counter\",\"delta\":111}",
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewMetricsStorage(t)
			if tt.mockGauge.needed {
				mockStorage.On("GetGauge", mock.Anything, tt.mockGauge.name).Return(tt.mockGauge.value, tt.mockGauge.err)
			}
			if tt.mockCounter.needed {
				mockStorage.On("GetCounter", mock.Anything, tt.mockCounter.name).Return(tt.mockCounter.value, tt.mockCounter.err)
			}

			r := chi.NewRouter()
			r.Post("/value", GetBody(mockStorage))

			req := httptest.NewRequest(tt.method, "/value", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			if res.StatusCode == http.StatusOK {
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.want.body, string(data))
			}

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			defer res.Body.Close()
		})
	}

}
