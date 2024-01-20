package main

import (
	"fmt"
	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/server"
	"github.com/vindosVP/metrics/pkg/logger"
	"log"
)

func main() {
	log.Print("Starting metrics server")
	cfg := config.NewServerConfig()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	err = server.Run(cfg)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to start server: %v", err))
	}
}
