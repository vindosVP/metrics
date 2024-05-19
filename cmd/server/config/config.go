package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v10"
)

type ServerConfig struct {
	RunAddr         string
	LogLevel        string
	FileStoragePath string
	DatabaseDNS     string
	Key             string
	StoreInterval   time.Duration
	Restore         bool
	CryptoKeyFile   string
}

type tempConfig struct {
	RunAddr         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDNS     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	CryptoKeyFile   string `env:"CRYPTO_KEY"`
	StoreInterval   int
	Restore         bool
}

func NewServerConfig() *ServerConfig {

	flagConfig := &tempConfig{}
	flag.StringVar(&flagConfig.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagConfig.LogLevel, "l", "info", "log level")
	flag.StringVar(&flagConfig.FileStoragePath, "f", "./tmp/metrics-db.json", "file storage path")
	flag.IntVar(&flagConfig.StoreInterval, "i", 300, "store interval")
	flag.BoolVar(&flagConfig.Restore, "r", true, "restore from dump file")
	flag.StringVar(&flagConfig.DatabaseDNS, "d", "", "database dns")
	flag.StringVar(&flagConfig.Key, "k", "", "hash key")
	flag.StringVar(&flagConfig.CryptoKeyFile, "crypto-key", "./keys/key.rsa", "crypto key")
	flag.Parse()

	envConfig := &tempConfig{}
	if err := env.Parse(envConfig); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	tempCfg := &tempConfig{}
	tempCfg.Restore = flagConfig.Restore
	tempCfg.StoreInterval = flagConfig.StoreInterval
	envRestore, ok := os.LookupEnv("RESTORE")
	if ok {
		restore, err := strconv.ParseBool(envRestore)
		if err != nil {
			log.Fatalf("Failed to parse env RESTORE value: %v", err)
		}
		tempCfg.Restore = restore
	}
	envStoreInterval, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		storeInterval, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			log.Fatalf("Failed to parse env STORE_INTERVAL value: %v", err)
		}
		tempCfg.StoreInterval = storeInterval
	}

	tempCfg.RunAddr = envConfig.RunAddr
	tempCfg.LogLevel = envConfig.LogLevel
	tempCfg.StoreInterval = envConfig.StoreInterval
	tempCfg.FileStoragePath = envConfig.FileStoragePath
	tempCfg.DatabaseDNS = envConfig.DatabaseDNS
	tempCfg.Key = envConfig.Key
	tempCfg.CryptoKeyFile = envConfig.CryptoKeyFile
	if tempCfg.Key == "" {
		tempCfg.Key = flagConfig.Key
	}
	if tempCfg.DatabaseDNS == "" {
		tempCfg.DatabaseDNS = flagConfig.DatabaseDNS
	}
	if tempCfg.RunAddr == "" {
		tempCfg.RunAddr = flagConfig.RunAddr
	}
	if tempCfg.LogLevel == "" {
		tempCfg.LogLevel = flagConfig.LogLevel
	}
	if tempCfg.FileStoragePath == "" {
		tempCfg.FileStoragePath = flagConfig.FileStoragePath
	}
	if tempCfg.CryptoKeyFile == "" {
		tempCfg.CryptoKeyFile = flagConfig.CryptoKeyFile
	}

	return &ServerConfig{
		RunAddr:         tempCfg.RunAddr,
		LogLevel:        tempCfg.LogLevel,
		FileStoragePath: tempCfg.FileStoragePath,
		Restore:         tempCfg.Restore,
		StoreInterval:   time.Duration(tempCfg.StoreInterval),
		DatabaseDNS:     tempCfg.DatabaseDNS,
		Key:             tempCfg.Key,
		CryptoKeyFile:   tempCfg.CryptoKeyFile,
	}
}
