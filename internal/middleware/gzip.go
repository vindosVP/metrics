package middleware

import (
	"compress/gzip"
	"github.com/vindosVP/metrics/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
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
