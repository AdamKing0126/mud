package sql_database

import (
	"database/sql"
	"mud/player"
)

// Database abstraction layer
type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// SQLite database implementation
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

func GetPlayer(db Database, id int) (*player.Player, error) {
	row := db.Query("SELECT id, name, location, health FROM players WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return nil, err
	}
	defer row.Close()

	player := &player.Player{}
	if err := row.Scan(&player.ID, &player.Name, &player.Location, &player.Health); err != nil {
		return nil, err
	}

	return player, nil
}
