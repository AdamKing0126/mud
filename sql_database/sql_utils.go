package sql_database

import (
	"database/sql"
)

type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SQLiteDatabase struct {
	db *sql.DB
}

func NewSQLiteDatabase(databaseFile string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return nil, err
	}
	return &SQLiteDatabase{db: db}, nil
}

func (db *SQLiteDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

func (db *SQLiteDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}
