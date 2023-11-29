package sql_database

import (
	"database/sql"
	"mud/players"
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

func GetPlayer(db Database, id int) (*players.Player, error) {
	row, err := db.Query("SELECT id, name, location, health FROM players WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	player := &players.Player{}
	if err := row.Scan(&player.UUID, &player.Name, &player.Room, &player.Area, &player.Health); err != nil {
		return nil, err
	}

	return player, nil
}
