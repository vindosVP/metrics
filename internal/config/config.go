package config

import (
	"flag"
)

type ServerConfig struct {
	RunAddr string
}

type AgentConfig struct {
	ServerAddr     string
	PollInterval   int
	ReportInterval int
}

func NewAgentConfig() *AgentConfig {
	config := &AgentConfig{}
	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.IntVar(&config.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&config.ReportInterval, "r", 10, "report interval")
	flag.Parse()
	return config
}

func NewServerConfig() *ServerConfig {
	config := &ServerConfig{}
	flag.StringVar(&config.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	return config
}
