package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateBatch(t *testing.T) {

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "ok",
			method: http.MethodPost,
			body:   "[{\"id\": \"test\",\"type\": \"gauge\",\"value\": 20},{\"id\": \"test\",\"type\": \"counter\",\"delta\": 10}]",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "wrong method",
			method: http.MethodGet,
			body:   "[{\"id\": \"test\",\"type\": \"gauge\",\"value\": 20},{\"id\": \"test\",\"type\": \"counter\",\"delta\": 10}]",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name:   "no gauge value",
			method: http.MethodPost,
			body:   "[{\"id\": \"test\",\"type\": \"gauge\",\"delta\": 20},{\"id\": \"test\",\"type\": \"counter\",\"delta\": 10}]",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "no counter delta",
			method: http.MethodPost,
			body:   "[{\"id\": \"test\",\"type\": \"gauge\",\"value\": 20},{\"id\": \"test\",\"type\": \"counter\",\"value\": 10}]",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "wrong counter delta",
			method: http.MethodPost,
			body:   "[{\"id\": \"test\",\"type\": \"gauge\",\"value\": 20},{\"id\": \"test\",\"type\": \"counter\",\"delta\": 10.3}]",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cRepo := repos.NewCounterRepo()
			gRepo := repos.NewGaugeRepo()
			storage := memstorage.New(gRepo, cRepo)

			r := chi.NewRouter()
			r.Post("/updates", UpdateBatch(storage))

			req := httptest.NewRequest(tt.method, "/updates", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			defer res.Body.Close()
		})
	}
}
