// Package dbstorage is a metrics storage working with postgres to store metrics.
package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"

	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/storage"
	"github.com/vindosVP/metrics/pkg/logger"
)

var retryDelays = map[uint]time.Duration{
	0: 1 * time.Second,
	1: 3 * time.Second,
	2: 5 * time.Second,
}

// Storage consists of postgres connection pool.
type Storage struct {
	db *pgxpool.Pool
}

// New creates the Storage.
func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		db: pool,
	}
}

// InsertBatch method saves provided metrics values to the database.
func (s *Storage) InsertBatch(ctx context.Context, batch []*models.Metrics) error {
	return retry.Do(func() error {
		b := batch
		tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
		defer tx.Rollback(ctx)
		if err != nil {
			return err
		}
		gaugeQuery := "insert into gauges (id, value) values ($1, $2) on conflict (id) do update set value = $2"
		counterQuery := "insert into counters as t (id, value) values ($1, $2) on conflict (id) do update set value = t.value + $2"
		for _, metric := range b {
			switch metric.MType {
			case models.Counter:
				val := *metric.Delta
				if _, err := tx.Exec(ctx, counterQuery, metric.ID, val); err != nil {
					return err
				}
			case models.Gauge:
				val := *metric.Value
				if _, err := tx.Exec(ctx, gaugeQuery, metric.ID, val); err != nil {
					return err
				}
			}
		}
		err = tx.Commit(ctx)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Error commiting transaction: %v", err))
		}
		return err
	}, retryOpts()...)
}

// UpdateGauge method updates gauge metric value.
// new value replaces the old one.
func (s *Storage) UpdateGauge(ctx context.Context, name string, v float64) (float64, error) {
	return retry.DoWithData(func() (float64, error) {
		query := "insert into gauges (id, value) values ($1, $2) on conflict (id) do update set value = $2"
		_, err := s.db.Exec(ctx, query, name, v)
		if err != nil {
			return 0, err
		}
		return v, nil
	}, retryOpts()...)
}

// UpdateCounter method updates counter metric value.
// new value adds to the old one.
func (s *Storage) UpdateCounter(ctx context.Context, name string, v int64) (int64, error) {
	return retry.DoWithData(func() (int64, error) {
		query := "insert into counters as t (id, value) values ($1, $2) on conflict (id) do update set value = t.value + $2"
		_, err := s.db.Exec(ctx, query, name, v)
		if err != nil {
			return 0, err
		}
		return v, nil
	}, retryOpts()...)
}

// GetGauge method returns value of gauge metric
func (s *Storage) GetGauge(ctx context.Context, name string) (float64, error) {
	return retry.DoWithData(func() (float64, error) {
		query := "select value from gauges where id = $1"
		row := s.db.QueryRow(ctx, query, name)
		var value float64
		err := row.Scan(&value)
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrMetricNotRegistered
		}
		if err != nil {
			return 0, err
		}
		return value, nil
	}, retryOpts()...)
}

// GetCounter method returns value of counter metric
func (s *Storage) GetCounter(ctx context.Context, name string) (int64, error) {
	return retry.DoWithData(func() (int64, error) {
		query := "select value from counters where id = $1"
		row := s.db.QueryRow(ctx, query, name)
		var value int64
		err := row.Scan(&value)
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrMetricNotRegistered
		}
		if err != nil {
			return 0, err
		}
		return value, nil
	}, retryOpts()...)
}

// GetAllGauge method returns values of all collected gauge metrics
func (s *Storage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return retry.DoWithData(func() (map[string]float64, error) {
		rows, err := s.db.Query(ctx, "select id, value from gauges order by id")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		res := make(map[string]float64)
		for rows.Next() {
			var id string
			var value float64
			err := rows.Scan(&id, &value)
			if err != nil {
				return nil, err
			}
			res[id] = value
		}

		if rows.Err() != nil {
			return nil, err
		}
		return res, nil
	}, retryOpts()...)
}

// GetAllCounter method returns values of all collected counter metrics
func (s *Storage) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	return retry.DoWithData(func() (map[string]int64, error) {
		rows, err := s.db.Query(ctx, "select id, value from counters order by id")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		res := make(map[string]int64)
		for rows.Next() {
			var id string
			var value int64
			err := rows.Scan(&id, &value)
			if err != nil {
				return nil, err
			}
			res[id] = value
		}

		if rows.Err() != nil {
			return nil, err
		}
		return res, nil
	}, retryOpts()...)
}

// SetCounter method sets counter metric value.
// new value replaces the old one.
func (s *Storage) SetCounter(ctx context.Context, name string, v int64) (int64, error) {
	return retry.DoWithData(func() (int64, error) {
		query := "insert into counters (id, value) values ($1, $2) on conflict (id) do update set value = $2"
		_, err := s.db.Exec(ctx, query, name, v)
		if err != nil {
			return 0, err
		}
		return v, nil
	}, retryOpts()...)
}

func retryOpts() []retry.Option {
	return []retry.Option{
		retry.RetryIf(func(err error) bool {
			return pgerrcode.IsConnectionException(pgErrCode(err)) || errors.Is(err, syscall.ECONNREFUSED)
		}),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			delay := retryDelays[n]
			return delay
		}),
		retry.OnRetry(func(n uint, err error) {
			logger.Log.Info(fmt.Sprintf("Failed to connect to database, retrying in %s", retryDelays[n]))
		}),
		retry.Attempts(4),
		retry.LastErrorOnly(true),
	}
}

func pgErrCode(err error) string {
	if e, ok := err.(*pq.Error); ok {
		return string(e.Code)
	}

	return ""
}
