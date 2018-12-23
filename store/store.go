package store

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

// InitMainDB - initialize the DB connection string
func InitMainDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("Something went wrong in MainDB init")
	}
	if e := InitFrontends(db); e != nil {
		return nil, e
	}
	if e := InitLocalRoutes(db); e != nil {
		return nil, e
	}
	if e := InitBackends(db); e != nil {
		return nil, e
	}
	if e := InitTargets(db); e != nil {
		return nil, e
	}
	if e := InitTargetLists(db); e != nil {
		return nil, e
	}
	return db, nil
}
