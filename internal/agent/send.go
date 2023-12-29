package agent

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func send(s *storage, wg *sync.WaitGroup) {
	for true {
		s.mu.Lock()
		for name, val := range s.gaugeMetrics {
			sendGaugeMetric(name, val)
		}
		for name, val := range s.counterMetrics {
			sendCounterMetric(name, val)
		}
		s.mu.Unlock()
		time.Sleep(reportInterval * time.Second)
	}
	wg.Done()
}

func sendGaugeMetric(name string, value float64) {
	URL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverHost, gauge, name, value)
	res, err := http.Post(URL, "text/plain", nil)
	if err != nil {
		log.Print(fmt.Sprintf("Failed to send metric %s: %v", name, err))
		return
	}
	if res.StatusCode != 200 {
		log.Print(fmt.Sprintf("Server respond %d on metric %s", res.StatusCode, name))
		return
	}
	log.Print(fmt.Sprintf("Metric %s sent", name))
}

func sendCounterMetric(name string, value int64) {
	URL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverHost, counter, name, value)
	res, err := http.Post(URL, "text/plain", nil)
	if err != nil {
		log.Print(fmt.Sprintf("Failed to send metric %s: %v", name, err))
		return
	}
	if res.StatusCode != 200 {
		log.Print(fmt.Sprintf("Server respond %d on metric %s", res.StatusCode, name))
		return
	}
	log.Print(fmt.Sprintf("Metric %s sent", name))
}
