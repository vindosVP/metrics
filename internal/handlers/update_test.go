package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
)

func ExampleUpdate() {
	// create new router
	r := chi.NewRouter()

	// init server config
	cfg := config.NewServerConfig()

	// create storage
	s := memstorage.New(repos.NewGaugeRepo(), repos.NewCounterRepo())

	// register handler
	r.Post("/update/{type}/{name}/{value}", UpdateBody(s))

	// start server
	log.Fatal(http.ListenAndServe(cfg.RunAddr, r))
}

func TestUpdate(t *testing.T) {

	type want struct {
		contentType string
		code        int
	}

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "gauge ok",
			url:    "/update/gauge/Alloc/12.5",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "counter ok",
			url:    "/update/counter/PollCount/12",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "wrong method",
			url:    "/update/counter/PollCount/12",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name:   "gauge bad request",
			url:    "/update/gauge/Alloc/true",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "counter bad request",
			url:    "/update/counter/PollCount/11.3",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "invalid type",
			url:    "/update/test/PollCount/11.3",
			method: http.MethodPost,
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
			r.Post("/update/{type}/{name}/{value}", Update(storage))

			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			defer res.Body.Close()
		})
	}

}
