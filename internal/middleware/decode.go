package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/vindosVP/metrics/pkg/encryption"
	"github.com/vindosVP/metrics/pkg/logger"
)

type decoder struct {
	cryptoKey *rsa.PrivateKey
}

func newDecoder(key *rsa.PrivateKey) *decoder {
	return &decoder{cryptoKey: key}
}

func Decode(key *rsa.PrivateKey) func(next http.Handler) http.Handler {
	d := newDecoder(key)
	return d.DecodeHandler
}

func (d *decoder) DecodeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b bytes.Buffer
		_, err := io.Copy(&b, r.Body)
		if err != nil {
			logger.Log.Error("Failed to read request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dec, err := encryption.Decrypt(d.cryptoKey, b.Bytes())
		if err != nil {
			logger.Log.Error("Failed to decrypt request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(dec))

		next.ServeHTTP(w, r)
	})
}
