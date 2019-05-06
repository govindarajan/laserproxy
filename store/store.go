package store

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

// InitMainDB - initialize the DB connection string
func InitMainDB(db *sql.DB) error {
	if db == nil {
		return errors.New("Something went wrong in MainDB init")
	}
	if e := InitFrontend(db); e != nil {
		return e
	}
	if e := InitLocalRoute(db); e != nil {
		return e
	}
	if e := InitBackend(db); e != nil {
		return e
	}
	if e := InitTargets(db); e != nil {
		return e
	}
	if e := InitTargetLists(db); e != nil {
		return e
	}
	return nil
}

// GetConnection used to get main DB connection.
func GetConnection() (*sql.DB, error) {
	return sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
}
