package dbstorage

import (
	"context"
	"database/sql"
	"github.com/vindosVP/metrics/internal/models"
)

type Counter interface {
	Insert(ctx context.Context, name string, v int64) (int64, error)
	Update(ctx context.Context, name string, v int64) (int64, error)
	Exists(ctx context.Context, name string) (bool, error)
	Get(ctx context.Context, name string) (int64, error)
	GetAll(ctx context.Context) (map[string]int64, error)
}

type Gauge interface {
	Insert(ctx context.Context, name string, v float64) (float64, error)
	Update(ctx context.Context, name string, v float64) (float64, error)
	Exists(ctx context.Context, name string) (bool, error)
	Get(ctx context.Context, name string) (float64, error)
	GetAll(ctx context.Context) (map[string]float64, error)
}

type Storage struct {
	db *sql.DB
	cr Counter
	gr Gauge
}

func New(cRepo Counter, gRepo Gauge, db *sql.DB) *Storage {
	return &Storage{
		db: db,
		cr: cRepo,
		gr: gRepo,
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
	cstmt, err := tx.PrepareContext(ctx, "insert into counters (id, value) values ($1, $2) on conflict (id) do update set value = $2")
	if err != nil {
		return err
	}
	for _, metric := range batch {
		switch metric.MType {
		case models.Counter:
			exists, err := s.cr.Exists(ctx, metric.ID)
			if err != nil {
				return err
			}
			val := *metric.Delta
			if exists {
				cval, err := s.cr.Get(ctx, metric.ID)
				if err != nil {
					return err
				}
				val += cval
			}
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
	exists, err := s.gr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if exists {
		return s.gr.Update(ctx, name, v)
	}
	return s.gr.Insert(ctx, name, v)
}

func (s *Storage) UpdateCounter(ctx context.Context, name string, v int64) (int64, error) {
	exists, err := s.cr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return s.cr.Insert(ctx, name, v)
	}
	cval, err := s.cr.Get(ctx, name)
	if err != nil {
		return 0, err
	}
	return s.cr.Update(ctx, name, cval+v)
}

func (s *Storage) GetGauge(ctx context.Context, name string) (float64, error) {
	exists, err := s.gr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrMetricNotRegistered
	}
	return s.gr.Get(ctx, name)
}

func (s *Storage) GetCounter(ctx context.Context, name string) (int64, error) {
	exists, err := s.cr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrMetricNotRegistered
	}
	return s.cr.Get(ctx, name)
}

func (s *Storage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return s.gr.GetAll(ctx)
}

func (s *Storage) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	return s.cr.GetAll(ctx)
}

func (s *Storage) SetCounter(ctx context.Context, name string, v int64) (int64, error) {
	exists, err := s.cr.Exists(ctx, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return s.cr.Insert(ctx, name, v)
	}
	return s.cr.Update(ctx, name, v)
}
