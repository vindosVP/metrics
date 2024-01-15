package config

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"log"
)

type ServerConfig struct {
	RunAddr string `env:"ADDRESS"`
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
