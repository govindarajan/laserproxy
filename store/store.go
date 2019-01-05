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
	if e := InitFrontends(db); e != nil {
		return e
	}
	if e := InitLocalRoute(db); e != nil {
		return e
	}
	if e := InitBackends(db); e != nil {
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
	return sql.Open("sqlite3", ":memory:")
}
