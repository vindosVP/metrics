// Package agent is a simple package to collect system metrics and sent them to the server.
package agent

import (
	"crypto/rsa"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/agent/collector"
	"github.com/vindosVP/metrics/internal/agent/sender"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"github.com/vindosVP/metrics/pkg/encryption"
	"github.com/vindosVP/metrics/pkg/logger"
)

// Run starts the agent
func Run(cfg *config.AgentConfig) error {

	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)

	c := collector.New(cfg.PollInterval, storage)
	var key *rsa.PublicKey = nil
	if cfg.CryptoKeyFile != "" {
		k, err := encryption.PublicKeyFromFile(cfg.CryptoKeyFile)
		if err != nil {
			logger.Log.Fatal("failed to get encryption key", zap.Error(err))
		}
		key = k
	}
	s := sender.New(cfg, storage, key)

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sig
		logger.Log.Info("Got stop signal, stopping")
		c.Stop()
		s.Stop()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go c.Run(&wg)
	go s.Run(&wg)
	wg.Wait()

	logger.Log.Info("Stopped successfully")

	return nil
}
