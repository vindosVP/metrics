package main

import (
	"fmt"
	"github.com/vindosVP/metrics/internal/server"
	"log"
)

func main() {
	log.Print("Starting metrics server")
	err := server.Run()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to start server: %v", err))
	}
}
