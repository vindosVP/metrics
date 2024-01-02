package agent

import (
	"flag"
	"time"
)

type Config struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerAddr     string
}

func NewConfig() *Config {
	config := &Config{}
	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.DurationVar(&config.PollInterval, "p", 2, "metrics poll interval")
	flag.DurationVar(&config.ReportInterval, "r", 10, "report interval")
	flag.Parse()
	return config
}
