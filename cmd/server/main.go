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
	err = server.Run(cfg)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func printBuildInfo() {
	fmt.Println(fmt.Sprintf("Build version: %s", buildVersion))
	fmt.Println(fmt.Sprintf("Build date: %s", buildDate))
	fmt.Println(fmt.Sprintf("Build commit: %s", buildCommit))
}
