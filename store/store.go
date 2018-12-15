package store

import (
	"database/sql"
	"fmt"

	"github.com/govindarajan/laserproxy/logger"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB - initialize the DB connection string and
// also create required tables
func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		logger.LogError("Exception at SQLite action: %+v" + err.Error())
		return nil
	}
	if db == nil {
		panic("DB nil")
	}
	return db
}

// Write or update rows with parameterized query String
func Write(db *sql.DB, StoreSQLStr string) error {
	_, err := db.Exec(StoreSQLStr)
	if err != nil {
		logger.LogError("Exception at SQLite action: %+v" + err.Error())
		return err
	}
	return nil
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
		fmt.Println("SQLite: Exception at et columns", err)
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
			fmt.Println("Failed to scan row", err)
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
