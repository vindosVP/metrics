package storage

import "errors"

var (
	ErrMetricNotRegistered = errors.New("metric with this name not registered")
)
