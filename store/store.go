package store

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB - initialize the DB connection string
func InitDB(db *sql.DB, dName string) error {
	if db == nil {
		return errors.New("Something went wrong in DB init " + dName)
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
	if e := InitTargetLists(db, dName); e != nil {
		return e
	}
	return nil
}

// GetConnection used to get main DB connection.
func GetConnection() (*sql.DB, error) {
	return sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
}

var fileDB = "/var/lib/laserproxy.db"

func GetFileDB() string {
	return fileDB
}

func GetFileDBConn() (*sql.DB, error) {
	return sql.Open("sqlite3", "file:"+fileDB+"?mode=rwc&cache=shared")
}

func MoveConfigToDisk(db *sql.DB) error {
	//.
	sql := "attach database '" + GetFileDB() + "' as disk"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy Frontend
	sql = "DELETE FROM disk.Frontend; INSERT INTO disk.Frontend SELECT * FROM Frontend;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy Backend
	sql = "DELETE FROM disk.Backend; INSERT INTO disk.Backend SELECT * FROM Backend;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy Targets
	sql = "DELETE FROM disk.Targets; INSERT INTO disk.Targets SELECT * FROM Targets;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy TargetLists
	sql = "DELETE FROM disk.TargetLists; INSERT INTO disk.TargetLists SELECT * FROM TargetLists;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	sql = "detach database disk;"
	_, err = db.Exec(sql)

	return err
}

func LoadConfigFromDisk(db *sql.DB) error {
	//.
	sql := "attach database '" + GetFileDB() + "' as disk"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy Frontend
	sql = "DELETE FROM Frontend; INSERT INTO Frontend SELECT * FROM disk.Frontend;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy Backend
	sql = "DELETE FROM Backend; INSERT INTO Backend SELECT * FROM disk.Backend;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy Targets
	sql = "DELETE FROM Targets; INSERT INTO Targets SELECT * FROM disk.Targets;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Copy TargetLists
	sql = "DELETE FROM TargetLists; INSERT INTO TargetLists SELECT * FROM disk.TargetLists;"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	sql = "detach database disk;"
	_, err = db.Exec(sql)

	return err
}
