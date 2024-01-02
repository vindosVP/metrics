package agent

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func send(s *storage, cfg *Config, wg *sync.WaitGroup) {
	for {
		s.mu.Lock()
		for name, val := range s.gaugeMetrics {
			sendGaugeMetric(name, val, cfg.ServerAddr)
		}
		for name, val := range s.counterMetrics {
			sendCounterMetric(name, val, cfg.ServerAddr)
		}
		s.mu.Unlock()
		time.Sleep(cfg.ReportInterval * time.Second)
	}
}

func sendGaugeMetric(name string, value float64, serverAddr string) {
	URL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverAddr, gauge, name, value)
	res, err := http.Post(URL, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to send metric %s: %v", name, err)
		return
	}
	if res.StatusCode != 200 {
		log.Printf("Server respond %d on metric %s", res.StatusCode, name)
		return
	}
	log.Printf("Metric %s sent", name)
	defer res.Body.Close()
}

func sendCounterMetric(name string, value int64, serverAddr string) {
	URL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverAddr, counter, name, value)
	res, err := http.Post(URL, "text/plain", nil)
	if err != nil {
		log.Printf("Failed to send metric %s: %v", name, err)
		return
	}
	if res.StatusCode != 200 {
		log.Printf("Server respond %d on metric %s", res.StatusCode, name)
		return
	}
	log.Printf("Metric %s sent", name)
	defer res.Body.Close()
}
