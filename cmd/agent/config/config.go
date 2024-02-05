package config

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"log"
	"time"
)

type AgentConfig struct {
	ServerAddr     string
	PollInterval   time.Duration
	ReportInterval time.Duration
	LogLevel       string
}

type tempConfig struct {
	ServerAddr     string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
}

func NewAgentConfig() *AgentConfig {

	flagConfig := &tempConfig{}
	flag.StringVar(&flagConfig.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.IntVar(&flagConfig.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&flagConfig.ReportInterval, "r", 10, "report interval")
	flag.StringVar(&flagConfig.LogLevel, "l", "info", "log level")
	flag.Parse()

	envConfig := &tempConfig{}
	if err := env.Parse(envConfig); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	tempCfg := &tempConfig{}
	tempCfg.ServerAddr = envConfig.ServerAddr
	tempCfg.PollInterval = envConfig.PollInterval
	tempCfg.ReportInterval = envConfig.ReportInterval
	tempCfg.LogLevel = envConfig.LogLevel
	if tempCfg.ServerAddr == "" {
		tempCfg.ServerAddr = flagConfig.ServerAddr
	}
	if tempCfg.PollInterval == 0 {
		tempCfg.PollInterval = flagConfig.PollInterval
	}
	if tempCfg.ReportInterval == 0 {
		tempCfg.ReportInterval = flagConfig.ReportInterval
	}
	if tempCfg.LogLevel == "" {
		tempCfg.LogLevel = flagConfig.LogLevel
	}

	config := &AgentConfig{
		ServerAddr:     tempCfg.ServerAddr,
		PollInterval:   time.Duration(tempCfg.PollInterval),
		ReportInterval: time.Duration(tempCfg.ReportInterval),
		LogLevel:       tempCfg.LogLevel,
	}

	return config
}
