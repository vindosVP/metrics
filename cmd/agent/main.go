package main

import (
	"fmt"
	"github.com/vindosVP/metrics/internal/agent"
	"log"
)

func main() {
	log.Print("Starting agent")
	err := agent.Run()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to start agent: %v", err))
	}
}
