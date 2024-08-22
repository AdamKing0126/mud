package database

import (
	"context"
)

type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Queryx(ctx context.Context, query string, args ...interface{}) (Rows, error)
	MapScan(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error)
	Prepare(ctx context.Context, query string) (Stmt, error)
	Begin(ctx context.Context) (Tx, error)
	Close() error
}

type Stmt interface {
	Close() error
	Exec(ctx context.Context, args ...interface{}) error
	Query(ctx context.Context, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, args ...interface{}) Row
}

type Tx interface {
	Commit() error
	Rollback() error
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Exec(ctx context.Context, query string, args ...interface{}) error
}

type Row interface {
	Scan(dest ...interface{}) error
	StructScan(dest interface{}) error
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	StructScan(dest interface{}) error
	Columns() ([]string, error)
	Close() error
	Err() error
}
