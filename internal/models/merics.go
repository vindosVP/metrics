// Package models consists of models of entities
package models

const (
	// Counter - counter metric type
	Counter = "counter"

	// Gauge - gauge metric type
	Gauge = "gauge"
)

// Metrics - structure of metric
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// MetricsDump - structure of metrics dump
type MetricsDump struct {
	Metrics []*Metrics `json:"metrics"`
}
