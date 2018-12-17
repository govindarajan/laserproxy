package store

import (
	"database/sql"
	"os"

	"github.com/govindarajan/laserproxy/logger"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB - initialize the DB connection string
func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		logger.LogError("Exception at SQLite action: %+v" + err.Error())
		return nil
	}
	if db == nil {
		panic("DB nil")
	}
	// create table if not exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		sqlStmt := `
		CREATE TABLE IF NOT EXISTS liveRequests (id INTEGER PRIMARY KEY AUTOINCREMENT, method TEXT, requestURI TEXT, sourceID INTEGER, targetID INTGER);
		CREATE TABLE IF NOT EXISTS HTTPTargets (id INTEGER PRIMARY KEY AUTOINCREMENT, host TEXT, ip TEXT, port INTEGER, weight INTEGER, status TEXT, maxRequests INTEGER);
		CREATE TABLE IF NOT EXISTS requestCounter (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, method TEXT, count INTEGER, statusCode INTEGER);
		CREATE TABLE IF NOT EXISTS liveRoutes (id INTEGER PRIMARY KEY AUTOINCREMENT, host TEXT, ip TEXT, port INTEGER);
		CREATE TABLE IF NOT EXISTS HTTPSourceRoutes (id INTEGER PRIMARY KEY AUTOINCREMENT, Interface TEXT, Type TEXT, healthCheck INTEGER, internalGateway TEXT, externalGateway TEXT, weight INTEGER, status TEXT, MaxRequests INTEGER);
		`
		_, err := db.Exec(sqlStmt)
		if err != nil {
			logger.LogError("Exception at SQLite action: %+v" + err.Error())
			return nil
		}
	}
	return db

}

// Write or update rows with parameterized query String
func Write(db *sql.DB, StoreSQLStr string) (bool, error) {
	_, err := db.Exec(StoreSQLStr)
	if err != nil {
		logger.LogError("Exception at SQLite action: %+v" + err.Error())
		return false, err
	}
	return true, nil
}

// Read data from the table with queryString
func Read(db *sql.DB, selectSQLStr string) [][]string {
	rows, err := db.Query(selectSQLStr)
	if err != nil {
		logger.LogError("SQLite: Exception at SQLite read: %+v" + err.Error())
		return nil
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		logger.LogError("SQLite: Exception at SQLite read rows: %+v" + err.Error())
		return nil

	}
	cols, err := rows.Columns()
	if err != nil {
		logger.LogError("SQLite: Exception at et columns: %+v" + err.Error())
		return nil
	}
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}
	values := [][]string{}
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			logger.LogError("SQLite: Failed to scan row: %+v" + err.Error())
			return nil
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}
		values = append(values, result)
	}
	return values
}
