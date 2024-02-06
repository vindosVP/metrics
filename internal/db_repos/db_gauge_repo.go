package db_repos

import (
	"context"
	"database/sql"
)

type GaugeRepo struct {
	db *sql.DB
}

func NewGaugeRepo(db *sql.DB) *GaugeRepo {
	return &GaugeRepo{db: db}
}

func (gr *GaugeRepo) Insert(ctx context.Context, name string, v float64) (float64, error) {
	_, err := gr.db.ExecContext(ctx, "insert into gauges (id, value) values ($1, $2)", name, v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (gr *GaugeRepo) Update(ctx context.Context, name string, v float64) (float64, error) {
	_, err := gr.db.ExecContext(ctx, "update gauges set value = $1 where id = $2", v, name)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (gr *GaugeRepo) Exists(ctx context.Context, name string) (bool, error) {
	row := gr.db.QueryRowContext(ctx, "select count(*) from gauges where id = $1", name)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (gr *GaugeRepo) Get(ctx context.Context, name string) (float64, error) {
	row := gr.db.QueryRowContext(ctx, "select value from gauges where id = $1", name)
	var value float64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (gr *GaugeRepo) GetAll(ctx context.Context) (map[string]float64, error) {
	rows, err := gr.db.QueryContext(ctx, "select id, value from gauges order by id")
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
