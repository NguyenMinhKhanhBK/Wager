package database

import (
	"context"
	"database/sql"
	"errors"
)

type DBManager interface {
	BeginTx() (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecWithContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryWithContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type database struct {
	db *sql.DB
}

func NewDB(db *sql.DB) DBManager {
	return &database{db: db}
}

func (d *database) BeginTx() (*sql.Tx, error) {
	return d.db.Begin()
}

func (d *database) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("invalid transaction")
	}
	return tx.Commit()
}

func (d *database) RollbackTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("invalid transaction")
	}
	return tx.Rollback()
}

func (d *database) ExecWithContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *database) QueryWithContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}
