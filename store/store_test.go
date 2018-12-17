package store_test

import (
	"testing"

	_ "github.com/govindarajan/laserproxy/store"
)

func TestAll(t *testing.T) {
	const dbpath = "lazer.db"

	db := InitDB(dbpath)
	defer db.Close()
	// CreateTable(db)

	var SQLstr string
	SQLstr = "insert into liveRequests (nil, 'get', 'exotel.com', 1, 2);"
	Write(db, SQLstr)

	var RSQLstr string
	RSQLstr = "select * from liveRequests;"
	readRow := Read(db, RSQLstr)
	t.Log(readRow)
}
