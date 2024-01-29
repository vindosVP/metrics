package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var compressionTypes = map[string]bool{
	"text/html":        true,
	"application/json": true,
	"html/text":        true,
}

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithCompression(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		needCompression := compressionTypes[r.Header.Get("Content-Type")] || compressionTypes[r.Header.Get("Accept")]

		if supportGzip && needCompression {
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()
			ow = compressWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}

			w.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = gzipReader
			defer gzipReader.Close()
		}

		h.ServeHTTP(ow, r)
	}
}
