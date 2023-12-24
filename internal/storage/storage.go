package storage

type MetricsStorage interface {
	UpdateGauge(name string, v float64) (float64, error)
	UpdateCounter(name string, v int64) (int64, error)
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}
