package main

import (
	"fmt"
	"log"

	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/server"
	"github.com/vindosVP/metrics/pkg/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()
	log.Print("Starting metrics server")
	cfg := config.NewServerConfig()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	logger.Log.Info("Starting server")
	s, err := server.New(cfg)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
	s.Run()
}

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
