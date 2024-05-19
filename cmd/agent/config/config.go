package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	Config         string
	ServerAddr     string
	LogLevel       string
	Key            string
	RateLimit      int
	PollInterval   time.Duration
	ReportInterval time.Duration
	CryptoKeyFile  string
}

type tempConfig struct {
	Config         string `json:"-"`
	ServerAddr     string `json:"address"`
	LogLevel       string `json:"log_level"`
	Key            string `json:"key"`
	RateLimit      int    `json:"rate_limit"`
	PollInterval   int    `json:"poll_interval"`
	ReportInterval int    `json:"report_interval"`
	CryptoKeyFile  string `json:"crypto_key_file"`
}

type configFullness struct {
	Config         bool
	ServerAddr     bool
	LogLevel       bool
	Key            bool
	RateLimit      bool
	PollInterval   bool
	ReportInterval bool
	CryptoKeyFile  bool
}

func NewAgentConfig() *AgentConfig {
	config := &AgentConfig{}
	full := &configFullness{}
	parseEnvs(config, full)
	parseFlags(config, full)
	if full.Config {
		parseJSON(config, full)
	}
	return config
}

func parseEnvs(config *AgentConfig, full *configFullness) {
	if val, ok := os.LookupEnv("CONFIG"); ok {
		config.Config = val
		full.Config = true
	}
	if val, ok := os.LookupEnv("ADDRESS"); ok {
		config.ServerAddr = val
		full.ServerAddr = true
	}
	if val, ok := os.LookupEnv("LOG_LEVEL"); ok {
		config.LogLevel = val
		full.LogLevel = true
	}
	if val, ok := os.LookupEnv("KEY"); ok {
		config.Key = val
		full.Key = true
	}
	if val, ok := os.LookupEnv("RATE_LIMIT"); ok {
		l, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Failed to parse env RATE_LIMIT value: %v", err)
		}
		config.RateLimit = l
		full.RateLimit = true
	}
	if val, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		p, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Failed to parse env POLL_INTERVAL value: %v", err)
		}
		config.PollInterval = time.Duration(p)
		full.PollInterval = true
	}
	if val, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		r, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Failed to parse env REPORT_INTERVAL value: %v", err)
		}
		config.ReportInterval = time.Duration(r)
		full.ReportInterval = true
	}
	if val, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		config.CryptoKeyFile = val
		full.CryptoKeyFile = true
	}
}

func parseFlags(config *AgentConfig, full *configFullness) {
	flagCfg := flags()
	if !full.Config && flagCfg.Config != "" {
		config.Config = flagCfg.Config
		full.Config = true
	}
	if !full.ServerAddr && flagCfg.ServerAddr != "" {
		config.ServerAddr = flagCfg.ServerAddr
		full.ServerAddr = true
	}
	if !full.LogLevel && flagCfg.LogLevel != "" {
		config.LogLevel = flagCfg.LogLevel
		full.LogLevel = true
	}
	if !full.Key && flagCfg.Key != "" {
		config.Key = flagCfg.Key
		full.Key = true
	}
	if !full.RateLimit && flagCfg.RateLimit != 0 {
		config.RateLimit = flagCfg.RateLimit
		full.RateLimit = true
	}
	if !full.PollInterval && flagCfg.PollInterval != 0 {
		config.PollInterval = time.Duration(flagCfg.PollInterval)
		full.PollInterval = true
	}
	if !full.ReportInterval {
		config.ReportInterval = time.Duration(flagCfg.ReportInterval)
		full.ReportInterval = true
	}
	if !full.CryptoKeyFile && flagCfg.CryptoKeyFile != "" {
		config.CryptoKeyFile = flagCfg.CryptoKeyFile
		full.CryptoKeyFile = true
	}
}

func parseJSON(config *AgentConfig, full *configFullness) {

	data, err := os.ReadFile(config.Config)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	JSONCfg := &tempConfig{}
	err = json.Unmarshal(data, &JSONCfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}
	if !full.ServerAddr && JSONCfg.ServerAddr != "" {
		config.ServerAddr = JSONCfg.ServerAddr
		full.ServerAddr = true
	}
	if !full.LogLevel && JSONCfg.LogLevel != "" {
		config.LogLevel = JSONCfg.LogLevel
		full.LogLevel = true
	}
	if !full.ReportInterval && JSONCfg.ReportInterval != 0 {
		config.ReportInterval = time.Duration(JSONCfg.ReportInterval)
		full.ReportInterval = true
	}
	if !full.PollInterval && JSONCfg.PollInterval != 0 {
		config.PollInterval = time.Duration(JSONCfg.PollInterval)
		full.PollInterval = true
	}
	if !full.Key && JSONCfg.Key != "" {
		config.Key = JSONCfg.Key
		full.Key = true
	}
	if !full.RateLimit {
		config.RateLimit = JSONCfg.RateLimit
		full.RateLimit = true
	}
	if !full.CryptoKeyFile && JSONCfg.CryptoKeyFile != "" {
		config.CryptoKeyFile = JSONCfg.CryptoKeyFile
		full.CryptoKeyFile = true
	}
}

func flags() *tempConfig {
	flagConfig := &tempConfig{}
	flag.StringVar(&flagConfig.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.IntVar(&flagConfig.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&flagConfig.ReportInterval, "r", 10, "report interval")
	flag.StringVar(&flagConfig.LogLevel, "lg", "info", "log level")
	flag.StringVar(&flagConfig.Key, "k", "", "secret key")
	flag.StringVar(&flagConfig.CryptoKeyFile, "crypto-key", "./keys/key.rsa.pub", "crypto key")
	flag.IntVar(&flagConfig.RateLimit, "l", 3, "rate limit")
	flag.Parse()
	return flagConfig
}
