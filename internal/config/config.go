package config

import (
	"flag"
)

type Config struct {
	RunAddr string
}

func New() *Config {
	config := &Config{}
	flag.StringVar(&config.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	return config
}
