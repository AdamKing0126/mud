package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Types

type SQLiteDB struct {
	db *sqlx.DB
}

type SQLiteStmt struct {
	stmt *sqlx.Stmt
}

type SQLiteRow struct {
	row *sqlx.Row
}

type SQLiteRows struct {
	rows *sqlx.Rows
}

type SQLiteTx struct {
	tx *sqlx.Tx
}

// SQLiteDB functions

func NewSQLiteDB(dataSourceName string) (*SQLiteDB, error) {
	db, err := sqlx.Connect("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *SQLiteDB) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	row := s.db.QueryRowxContext(ctx, query, args...)
	return &SQLiteRow{row: row}
}

func (s *SQLiteDB) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLiteRows{rows: rows}, nil
}

func (db *SQLiteDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.db.SelectContext(ctx, dest, query, args...)
}

func (db *SQLiteDB) Close() error {
	return db.db.Close()
}

func (db *SQLiteDB) Queryx(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := db.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLiteRows{rows: rows}, nil
}

func (db *SQLiteDB) MapScan(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := db.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	result := make(map[string]interface{})
	err = rows.MapScan(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *SQLiteDB) Prepare(ctx context.Context, query string) (Stmt, error) {
	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &SQLiteStmt{stmt: stmt}, nil
}

func (db *SQLiteDB) Begin(ctx context.Context) (Tx, error) {
	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &SQLiteTx{tx: tx}, nil
}

// SQLiteStmt functions

func (s *SQLiteStmt) Close() error {
	return s.stmt.Close()
}

func (s *SQLiteStmt) Exec(ctx context.Context, args ...interface{}) error {
	_, err := s.stmt.ExecContext(ctx, args...)
	return err
}

func (s *SQLiteStmt) Query(ctx context.Context, args ...interface{}) (Rows, error) {
	rows, err := s.stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return &SQLiteRows{rows: rows}, nil
}

func (s *SQLiteStmt) QueryRow(ctx context.Context, args ...interface{}) Row {
	row := s.stmt.QueryRowxContext(ctx, args...)
	return &SQLiteRow{row: row}
}

// SQLiteRow functions

func (r *SQLiteRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

func (r *SQLiteRow) StructScan(dest interface{}) error {
	return r.row.StructScan(dest)
}

// SQLiteRows functions

func (r *SQLiteRows) Next() bool {
	return r.rows.Next()
}

func (r *SQLiteRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *SQLiteRows) StructScan(dest interface{}) error {
	return r.rows.StructScan(dest)
}

func (r *SQLiteRows) Columns() ([]string, error) {
	return r.rows.Columns()
}

func (r *SQLiteRows) Close() error {
	return r.rows.Close()
}

func (r *SQLiteRows) Err() error {
	return r.rows.Err()
}

// SQLiteTx functions

func (tx *SQLiteTx) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := tx.tx.ExecContext(ctx, query, args...)
	return err
}

func (tx *SQLiteTx) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := tx.tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLiteRows{rows: rows}, nil
}

func (tx *SQLiteTx) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	row := tx.tx.QueryRowxContext(ctx, query, args...)
	return &SQLiteRow{row: row}
}

func (tx *SQLiteTx) Commit() error {
	return tx.tx.Commit()
}

func (tx *SQLiteTx) Rollback() error {
	return tx.tx.Rollback()
}
