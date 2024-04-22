package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/vindosVP/metrics/pkg/logger"
)

func Decompress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Log.Error("Failed to create gzip reader", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = gzipReader
			defer gzipReader.Close()
		}

		h.ServeHTTP(w, r)
	})
}
