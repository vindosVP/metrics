package main

import (
	"fmt"
	"log"

	"github.com/vindosVP/metrics/cmd/server/config"
	"github.com/vindosVP/metrics/internal/server"
	"github.com/vindosVP/metrics/pkg/logger"
)

func main() {
	log.Print("Starting metrics server")
	cfg := config.NewServerConfig()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	err = server.Run(cfg)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}

//func writeMemProfile() {
//	time.Sleep(60 * time.Second)
//	fmem, err := os.Create(`base.pprof`)
//	if err != nil {
//		panic(err)
//	}
//	defer fmem.Close()
//	runtime.GC() // получаем статистику по использованию памяти
//	if err := pprof.WriteHeapProfile(fmem); err != nil {
//		panic(err)
//	}
//}
