package store

import (
	"database/sql"
	"errors"
	"net"
)

type TargetLists struct {
	Hostname string
	IP       net.IP
	Score    int
}

func InitTargetLists(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS TargetLists (Hostname VARCHAR NOT NULL, IP VARCHAR NOT NULL,
		Score INT NOT NULL DEFAULT 0, KEY (hostname) );
		`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func WriteTargetLists(db *sql.DB, t *TargetLists) error {
	if t == nil {
		return errors.New("Empty target values are given")
	}
	stmt, err := db.Prepare("REPLACE INTO Targets (Hostname, IP, Score) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t.Hostname, t.IP, t.Score)
	return err
}
