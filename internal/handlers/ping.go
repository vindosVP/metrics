package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Ping(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := pool.Ping(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
