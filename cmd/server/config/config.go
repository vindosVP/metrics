package config

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"log"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	RunAddr         string
	LogLevel        string
	FileStoragePath string
	Restore         bool
	StoreInterval   time.Duration
}

type tempConfig struct {
	RunAddr         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool
	StoreInterval   int
}

func NewServerConfig() *ServerConfig {

	flagConfig := &tempConfig{}
	flag.StringVar(&flagConfig.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagConfig.LogLevel, "l", "info", "log level")
	flag.StringVar(&flagConfig.FileStoragePath, "f", "./tmp/metrics-db.json", "file storage path")
	flag.IntVar(&flagConfig.StoreInterval, "i", 300, "store interval")
	flag.BoolVar(&flagConfig.Restore, "r", true, "restore from dump file")
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
	if tempCfg.RunAddr == "" {
		tempCfg.RunAddr = flagConfig.RunAddr
	}
	if tempCfg.LogLevel == "" {
		tempCfg.LogLevel = flagConfig.LogLevel
	}
	if tempCfg.FileStoragePath == "" {
		tempCfg.FileStoragePath = flagConfig.FileStoragePath
	}

	return &ServerConfig{
		RunAddr:         tempCfg.RunAddr,
		LogLevel:        tempCfg.LogLevel,
		FileStoragePath: tempCfg.FileStoragePath,
		Restore:         tempCfg.Restore,
		StoreInterval:   time.Duration(tempCfg.StoreInterval),
	}
}
