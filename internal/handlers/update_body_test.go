package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
)

func TestUpdateBody(t *testing.T) {

	type want struct {
		code        int
		contentType string
		wantBody    bool
		body        string
	}

	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "gauge ok",
			method: http.MethodPost,
			body:   "{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":12.5}",
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				wantBody:    true,
				body:        "{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":12.5}",
			},
		},
		{
			name:   "counter ok",
			method: http.MethodPost,
			body:   "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":125}",
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				wantBody:    true,
				body:        "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":125}",
			},
		},
		{
			name:   "wrong method",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name:   "wrong type",
			method: http.MethodPost,
			body:   "{\"id\":\"Alloc\",\"type\":\"WrongType\",\"value\":12.5}",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "no id",
			method: http.MethodPost,
			body:   "{\"type\":\"gauge\",\"value\":12.5}",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "gauge wrong value",
			method: http.MethodPost,
			body:   "{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":true}",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "gauge no value",
			method: http.MethodPost,
			body:   "{\"id\":\"Alloc\",\"type\":\"gauge\"}",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "counter wrong value",
			method: http.MethodPost,
			body:   "{\"id\":\"PollCount\",\"type\":\"counter\",\"value\":true}",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "counter no value",
			method: http.MethodPost,
			body:   "{\"id\":\"PollCount\",\"type\":\"counter\"}",
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
			r.Post("/update", UpdateBody(storage))

			req := httptest.NewRequest(tt.method, "/update", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			if tt.want.wantBody {
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.want.body, string(data))
			}
			defer res.Body.Close()
		})
	}

}
