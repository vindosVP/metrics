package db_repos

import (
	"context"
	"database/sql"
)

type CounterRepo struct {
	db *sql.DB
}

func NewCounterRepo(db *sql.DB) *CounterRepo {
	return &CounterRepo{
		db: db,
	}
}

func (cr *CounterRepo) Insert(ctx context.Context, name string, v int64) (int64, error) {
	_, err := cr.db.ExecContext(ctx, "insert into counters (id, value) values ($1, $2)", name, v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (cr *CounterRepo) Update(ctx context.Context, name string, v int64) (int64, error) {
	_, err := cr.db.ExecContext(ctx, "update counters set value = $1 where id = $2", v, name)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (cr *CounterRepo) Exists(ctx context.Context, name string) (bool, error) {
	row := cr.db.QueryRowContext(ctx, "select count(*) from counters where id = $1", name)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (cr *CounterRepo) Get(ctx context.Context, name string) (int64, error) {
	row := cr.db.QueryRowContext(ctx, "select value from counters where id = $1", name)
	var value int64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (cr *CounterRepo) GetAll(ctx context.Context) (map[string]int64, error) {
	rows, err := cr.db.QueryContext(ctx, "select id, value from counters order by id")
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
