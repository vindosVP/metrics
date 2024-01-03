package agent

import (
	"github.com/vindosVP/metrics/internal/collector"
	"github.com/vindosVP/metrics/internal/config"
	"github.com/vindosVP/metrics/internal/repos"
	"github.com/vindosVP/metrics/internal/sender"
	"github.com/vindosVP/metrics/internal/storage/memstorage"
	"sync"
)

func Run() error {
	cfg := config.NewAgentConfig()

	cRepo := repos.NewCounterRepo()
	gRepo := repos.NewGaugeRepo()
	storage := memstorage.New(gRepo, cRepo)

	c := collector.New(cfg, storage)
	s := sender.New(cfg, storage)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go c.Run()
	go s.Run()
	wg.Wait()

	return nil
}
