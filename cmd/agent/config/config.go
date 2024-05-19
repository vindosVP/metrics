package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
)

type AgentConfig struct {
	ServerAddr     string
	LogLevel       string
	Key            string
	RateLimit      int
	PollInterval   time.Duration
	ReportInterval time.Duration
	CryptoKeyFile  string
}

type tempConfig struct {
	ServerAddr     string `env:"ADDRESS"`
	LogLevel       string `env:"LOG_LEVEL"`
	Key            string `env:"KEY"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKeyFile  string `env:"CRYPTO_KEY"`
}

func NewAgentConfig() *AgentConfig {

	flagConfig := &tempConfig{}
	flag.StringVar(&flagConfig.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.IntVar(&flagConfig.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&flagConfig.ReportInterval, "r", 10, "report interval")
	flag.StringVar(&flagConfig.LogLevel, "lg", "info", "log level")
	flag.StringVar(&flagConfig.Key, "k", "", "secret key")
	flag.StringVar(&flagConfig.CryptoKeyFile, "crypto-key", "./keys/key.rsa.pub", "crypto key")
	flag.IntVar(&flagConfig.RateLimit, "l", 3, "rate limit")
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
	tempCfg.Key = envConfig.Key
	tempCfg.RateLimit = envConfig.RateLimit
	tempCfg.CryptoKeyFile = envConfig.CryptoKeyFile
	if tempCfg.Key == "" {
		tempCfg.Key = flagConfig.Key
	}
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
	if tempCfg.RateLimit == 0 {
		tempCfg.RateLimit = flagConfig.RateLimit
	}
	if tempCfg.CryptoKeyFile == "" {
		tempCfg.CryptoKeyFile = flagConfig.CryptoKeyFile
	}

	config := &AgentConfig{
		ServerAddr:     tempCfg.ServerAddr,
		PollInterval:   time.Duration(tempCfg.PollInterval),
		ReportInterval: time.Duration(tempCfg.ReportInterval),
		LogLevel:       tempCfg.LogLevel,
		Key:            tempCfg.Key,
		RateLimit:      tempCfg.RateLimit,
		CryptoKeyFile:  tempCfg.CryptoKeyFile,
	}

	return config
}
