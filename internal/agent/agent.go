// Package agent is a simple package to collect system metrics and sent them to the server.
package agent

import (
	"sync"

	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/agent/collector"
	"github.com/vindosVP/metrics/internal/agent/sender"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
)

// Run starts the agent
func Run(cfg *config.AgentConfig) error {

	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)

	c := collector.New(cfg.PollInterval, storage)
	s := sender.New(cfg, storage)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go c.Run()
	go s.Run()
	wg.Wait()

	return nil
}
