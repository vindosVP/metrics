package config

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"log"
)

type ServerConfig struct {
	RunAddr string `env:"ADDRESS"`
}

type AgentConfig struct {
	ServerAddr     string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func NewAgentConfig() *AgentConfig {

	flagConfig := &AgentConfig{}
	flag.StringVar(&flagConfig.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.IntVar(&flagConfig.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&flagConfig.ReportInterval, "r", 10, "report interval")
	flag.Parse()

	envConfig := &AgentConfig{}
	if err := env.Parse(envConfig); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	cfg := &AgentConfig{}
	cfg.ServerAddr = envConfig.ServerAddr
	cfg.PollInterval = envConfig.PollInterval
	cfg.ReportInterval = envConfig.ReportInterval
	if cfg.ServerAddr == "" {
		cfg.ServerAddr = flagConfig.ServerAddr
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = flagConfig.PollInterval
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = flagConfig.ReportInterval
	}

	return cfg
}

func NewServerConfig() *ServerConfig {

	flagConfig := &ServerConfig{}
	flag.StringVar(&flagConfig.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	envConfig := &ServerConfig{}
	if err := env.Parse(envConfig); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	cfg := &ServerConfig{}
	cfg.RunAddr = envConfig.RunAddr
	if cfg.RunAddr == "" {
		cfg.RunAddr = flagConfig.RunAddr
	}

	return cfg
}
