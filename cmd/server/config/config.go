package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Config          string
	RunAddr         string
	LogLevel        string
	FileStoragePath string
	DatabaseDNS     string
	Key             string
	StoreInterval   time.Duration
	Restore         bool
	CryptoKeyFile   string
	TrustedSubnet   string
}

type tempConfig struct {
	Config          string
	RunAddr         string
	LogLevel        string
	FileStoragePath string
	DatabaseDNS     string
	Key             string
	StoreInterval   int
	Restore         bool
	CryptoKeyFile   string
	TrustedSubnet   string
}

type jsonConfig struct {
	RunAddr         string `json:"address"`
	LogLevel        string `json:"log_level"`
	FileStoragePath string `json:"store_file"`
	DatabaseDNS     string `json:"database_dsn"`
	Key             string `json:"key"`
	StoreInterval   int    `json:"store_interval"`
	Restore         bool   `json:"restore"`
	CryptoKeyFile   string `json:"crypto_key"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

type configFullness struct {
	Config          bool
	RunAddr         bool
	LogLevel        bool
	FileStoragePath bool
	DatabaseDNS     bool
	Key             bool
	CryptoKeyFile   bool
	StoreInterval   bool
	Restore         bool
	TrustedSubnet   bool
}

func NewServerConfig() *ServerConfig {
	config := &ServerConfig{}
	full := &configFullness{}
	parseEnvs(config, full)
	parseFlags(config, full)
	if full.Config {
		parseJSON(config, full)
	}
	return config
}

func parseFlags(config *ServerConfig, full *configFullness) {
	flagCfg := parseFlagConfig()
	if !full.Config && flagCfg.Config != "" {
		config.Config = flagCfg.Config
		full.Config = true
	}
	if !full.RunAddr && flagCfg.RunAddr != "" {
		config.RunAddr = flagCfg.RunAddr
		full.RunAddr = true
	}
	if !full.LogLevel && flagCfg.LogLevel != "" {
		config.LogLevel = flagCfg.LogLevel
		full.LogLevel = true
	}
	if !full.FileStoragePath && flagCfg.FileStoragePath != "" {
		config.FileStoragePath = flagCfg.FileStoragePath
		full.FileStoragePath = true
	}
	if !full.DatabaseDNS && flagCfg.DatabaseDNS != "" {
		config.DatabaseDNS = flagCfg.DatabaseDNS
		full.DatabaseDNS = true
	}
	if !full.Key && flagCfg.Key != "" {
		config.Key = flagCfg.Key
		full.Key = true
	}
	if !full.StoreInterval {
		config.StoreInterval = time.Duration(flagCfg.StoreInterval)
		full.StoreInterval = true
	}
	if !full.Restore {
		config.Restore = flagCfg.Restore
		full.Restore = true
	}
	if !full.TrustedSubnet && flagCfg.TrustedSubnet != "" {
		config.TrustedSubnet = flagCfg.TrustedSubnet
		full.TrustedSubnet = true
	}
	if !full.CryptoKeyFile && flagCfg.CryptoKeyFile != "" {
		config.CryptoKeyFile = flagCfg.CryptoKeyFile
		full.CryptoKeyFile = true
	}
}

func parseFlagConfig() *tempConfig {
	flagConfig := &tempConfig{}
	flag.StringVar(&flagConfig.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagConfig.LogLevel, "l", "info", "log level")
	flag.StringVar(&flagConfig.FileStoragePath, "f", "./tmp/metrics-db.json", "file storage path")
	flag.IntVar(&flagConfig.StoreInterval, "i", 300, "store interval")
	flag.BoolVar(&flagConfig.Restore, "r", true, "restore from dump file")
	flag.StringVar(&flagConfig.DatabaseDNS, "d", "", "database dns")
	flag.StringVar(&flagConfig.Key, "k", "", "hash key")
	flag.StringVar(&flagConfig.CryptoKeyFile, "crypto-key", "", "crypto key")
	flag.StringVar(&flagConfig.Config, "c", "", "json config file")
	flag.StringVar(&flagConfig.TrustedSubnet, "t", "", "trusted subnet")
	flag.Parse()
	return flagConfig
}

func parseEnvs(config *ServerConfig, full *configFullness) {
	if val, ok := os.LookupEnv("CONFIG"); ok {
		config.Config = val
		full.Config = true
	}
	if val, ok := os.LookupEnv("ADDRESS"); ok {
		config.RunAddr = val
		full.RunAddr = true
	}
	if val, ok := os.LookupEnv("LOG_LEVEL"); ok {
		config.LogLevel = val
		full.LogLevel = true
	}
	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoragePath = val
		full.FileStoragePath = true
	}
	if val, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DatabaseDNS = val
		full.DatabaseDNS = true
	}
	if val, ok := os.LookupEnv("KEY"); ok {
		config.Key = val
		full.Key = true
	}
	if val, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		config.CryptoKeyFile = val
		full.CryptoKeyFile = true
	}
	if val, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		config.TrustedSubnet = val
		full.TrustedSubnet = true
	}
	if val, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInterval, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Failed to parse env STORE_INTERVAL value: %v", err)
		}
		config.StoreInterval = time.Duration(storeInterval)
		full.StoreInterval = true
	}
	if val, ok := os.LookupEnv("RESTORE"); ok {
		restore, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("Failed to parse env RESTORE value: %v", err)
		}
		config.Restore = restore
		full.Restore = true
	}
}

func parseJSON(config *ServerConfig, full *configFullness) {

	data, err := os.ReadFile(config.Config)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	JSONCfg := &jsonConfig{}
	err = json.Unmarshal(data, &JSONCfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}
	if !full.RunAddr && JSONCfg.RunAddr != "" {
		config.RunAddr = JSONCfg.RunAddr
		full.RunAddr = true
	}
	if !full.LogLevel && JSONCfg.LogLevel != "" {
		config.LogLevel = JSONCfg.LogLevel
		full.LogLevel = true
	}
	if !full.FileStoragePath && JSONCfg.FileStoragePath != "" {
		config.FileStoragePath = JSONCfg.FileStoragePath
		full.FileStoragePath = true
	}
	if !full.DatabaseDNS && JSONCfg.DatabaseDNS != "" {
		config.DatabaseDNS = JSONCfg.DatabaseDNS
		full.DatabaseDNS = true
	}
	if !full.Key && JSONCfg.Key != "" {
		config.Key = JSONCfg.Key
		full.Key = true
	}
	if !full.StoreInterval {
		config.StoreInterval = time.Duration(JSONCfg.StoreInterval)
		full.StoreInterval = true
	}
	if !full.Restore {
		config.Restore = JSONCfg.Restore
		full.Restore = true
	}
	if !full.CryptoKeyFile && JSONCfg.CryptoKeyFile != "" {
		config.CryptoKeyFile = JSONCfg.CryptoKeyFile
		full.CryptoKeyFile = true
	}
	if !full.TrustedSubnet {
		config.TrustedSubnet = JSONCfg.TrustedSubnet
		full.TrustedSubnet = true
	}
}
