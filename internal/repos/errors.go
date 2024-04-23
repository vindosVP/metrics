package repos

import "errors"

var (

	// ErrMetricNotRegistered - represents that metric with provided name is not registered
	ErrMetricNotRegistered = errors.New("metric with this name not registered")
)
