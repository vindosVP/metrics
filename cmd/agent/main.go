package main

import (
	"fmt"
	"log"

	"github.com/vindosVP/metrics/cmd/agent/config"
	"github.com/vindosVP/metrics/internal/agent"
	"github.com/vindosVP/metrics/pkg/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()
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

func printBuildInfo() {
	fmt.Println(fmt.Sprintf("Build version: %s", buildVersion))
	fmt.Println(fmt.Sprintf("Build date: %s", buildDate))
	fmt.Println(fmt.Sprintf("Build commit: %s", buildCommit))
}
