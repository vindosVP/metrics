package main

import (
	"fmt"
	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/agent"
	"github.com/vindosVP/metrics/pkg/logger"
	"log"
)

func main() {
	log.Print("Starting agent")
	cfg := config.NewAgentConfig()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	err = agent.Run(cfg)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to start agent: %v", err))
	}
}
