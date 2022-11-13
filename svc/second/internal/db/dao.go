package db

import (
	"context"
	"fmt"
	"time"
)

type DAO struct{}

const createStmt = "INSERT INTO " + Table + " (`name`, `updated_at`) VALUES (?, ?);"

func (DAO) Create(ctx context.Context, name string) (*Model, error) {
	m := &Model{
		Name:      name,
		UpdatedAt: time.Now().UTC(),
	}
	res, err := db.ExecContext(ctx, createStmt, m.Name, m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("DAO#Create::ExecContext: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("DAO#Create::LastInsertId: %w", err)
	}
	m.ID = id
	return m, nil
}

const getByIDQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `id` = ? LIMIT 1;"

func (DAO) GetByID(ctx context.Context, id int64) (*Model, error) {
	rows, err := db.QueryContext(ctx, getByIDQuery, id)
	if err != nil {
		return nil, fmt.Errorf("DAO#GetByID::QueryContext id=%d: %w", id, err)
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		return nil, nil
	}
	var m Model
	if err := rows.Scan(m.ToPtrArgs()...); err != nil {
		return nil, fmt.Errorf("DAO#GetByID::Scan: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DAO#GetByID::rows.Err: %w", err)
	}
	return &m, nil
}
