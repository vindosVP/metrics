package middleware

import (
	"bytes"
	"github.com/vindosVP/metrics/pkg/logger"
	"github.com/vindosVP/metrics/pkg/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func ValidateHMAC(key string) func(next http.Handler) http.Handler {
	h := NewHasher(key)
	return h.ValidateHandler
}

func Sign(key string) func(next http.Handler) http.Handler {
	h := NewHasher(key)
	return h.SignHandler
}

type Hasher struct {
	key string
}

type responseSigner struct {
	w   http.ResponseWriter
	buf bytes.Buffer
	key string
}

func (r *responseSigner) Write(p []byte) (int, error) {
	r.buf.Write(p)
	return r.w.Write(p)
}

func (r *responseSigner) Header() http.Header {
	return r.w.Header()
}

func (r *responseSigner) WriteHeader(code int) {
	responseData := r.buf.Bytes()
	hash, err := utils.Sha256Hash(responseData, r.key)
	if err != nil {
		logger.Log.Error("Failed to compute hash", zap.Error(err))
	}
	if err == nil {
		r.w.Header().Set("HashSHA256", hash)
	}
	r.w.WriteHeader(code)
}

func NewHasher(key string) *Hasher {
	return &Hasher{key: key}
}

func (h *Hasher) SignHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.key != "" {
			signer := &responseSigner{
				w:   w,
				buf: bytes.Buffer{},
				key: h.key,
			}

			next.ServeHTTP(signer, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Hasher) ValidateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.key != "" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, r.Body)
			if err != nil {
				logger.Log.Error("Failed to read request body", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(&buf)
			hash, err := utils.Sha256Hash(buf.Bytes(), h.key)
			if err != nil {
				logger.Log.Error("Failed to compute hash", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			providedHash := r.Header.Get("HashSHA256")
			if providedHash != "" && providedHash != hash {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
