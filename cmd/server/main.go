package main

import (
	"fmt"
	"github.com/vindosVP/metrics/internal/config"
	"github.com/vindosVP/metrics/internal/server"
	"log"
)

func main() {
	log.Print("Starting metrics server")
	cfg := config.NewServerConfig()
	err := server.Run(cfg)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to start server: %v", err))
	}
}
