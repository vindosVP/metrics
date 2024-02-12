package dbstorage

import (
	"context"
	"database/sql"
	"github.com/vindosVP/metrics/internal/models"
	"github.com/vindosVP/metrics/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) InsertBatch(ctx context.Context, batch []*models.Metrics) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	gstmt, err := tx.PrepareContext(ctx, "insert into gauges (id, value) values ($1, $2) on conflict (id) do update set value = $2")
	if err != nil {
		return err
	}
	cstmt, err := tx.PrepareContext(ctx, "insert into counters as t (id, value) values ($1, $2) on conflict (id) do update set value = t.value + $2")
	if err != nil {
		return err
	}
	for _, metric := range batch {
		switch metric.MType {
		case models.Counter:
			val := *metric.Delta
			if _, err := cstmt.ExecContext(ctx, metric.ID, val); err != nil {
				return err
			}
		case models.Gauge:
			val := *metric.Value
			if _, err := gstmt.ExecContext(ctx, metric.ID, val); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *Storage) UpdateGauge(ctx context.Context, name string, v float64) (float64, error) {
	query := "insert into gauges (id, value) values ($1, $2) on conflict (id) do update set value = $2"
	_, err := s.db.ExecContext(ctx, query, name, v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (s *Storage) UpdateCounter(ctx context.Context, name string, v int64) (int64, error) {
	query := "insert into counters as t (id, value) values ($1, $2) on conflict (id) do update set value = t.value + $2"
	_, err := s.db.ExecContext(ctx, query, name, v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (s *Storage) GetGauge(ctx context.Context, name string) (float64, error) {
	query := "select value from gauges where id = $1"
	row := s.db.QueryRowContext(ctx, query, name)
	var value sql.NullFloat64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}
	if !value.Valid {
		return 0, storage.ErrMetricNotRegistered
	}
	return value.Float64, nil
}

func (s *Storage) GetCounter(ctx context.Context, name string) (int64, error) {
	query := "select value from counters where id = $1"
	row := s.db.QueryRowContext(ctx, query, name)
	var value sql.NullInt64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}
	if !value.Valid {
		return 0, storage.ErrMetricNotRegistered
	}
	return value.Int64, nil
}

func (s *Storage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	rows, err := s.db.QueryContext(ctx, "select id, value from gauges order by id")
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
}

func (s *Storage) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	rows, err := s.db.QueryContext(ctx, "select id, value from counters order by id")
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
}

func (s *Storage) SetCounter(ctx context.Context, name string, v int64) (int64, error) {
	query := "insert into counters (id, value) values ($1, $2) on conflict (id) do update set value = $2"
	_, err := s.db.ExecContext(ctx, query, name, v)
	if err != nil {
		return 0, err
	}
	return v, nil
}
