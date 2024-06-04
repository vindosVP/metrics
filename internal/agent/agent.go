// Package agent is a simple package to collect system metrics and sent them to the server.
package agent

import (
	"crypto/rsa"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/agent/collector"
	"github.com/vindosVP/metrics/internal/agent/sender"
	"github.com/vindosVP/metrics/internal/agent/senderRPC"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"github.com/vindosVP/metrics/pkg/encryption"
	"github.com/vindosVP/metrics/pkg/logger"
)

type Sender interface {
	Run(wg *sync.WaitGroup)
	Stop()
}

// Run starts the agent
func Run(cfg *config.AgentConfig) error {

	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)

	c := collector.New(cfg.PollInterval, storage)
	var s Sender
	if !cfg.UseRPC {
		var key *rsa.PublicKey = nil
		if cfg.CryptoKeyFile != "" {
			k, err := encryption.PublicKeyFromFile(cfg.CryptoKeyFile)
			if err != nil {
				logger.Log.Fatal("failed to get encryption key", zap.Error(err))
			}
			key = k
		}
		logger.Log.Info("Sending metrics using HTTP")
		s = sender.New(cfg, storage, key, GetLocalIP())
	} else {
		logger.Log.Info("Sending metrics using GRPC")
		s = senderRPC.New(cfg, storage)
	}

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

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
