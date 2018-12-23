package store

import (
	"database/sql"
	"errors"
	"net"
)

type LRStatus string

const (
	LRStOnline  LRStatus = "ONLINE"
	LRStShunned LRStatus = "SHUNNED"
)

type LocalRoutes struct {
	Interface     string
	IP            net.IP
	Gateway       net.IP
	CheckURL      *string
	CheckInterval int
	Weight        int
	Status        LRStatus
}

func InitLocalRoutes(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS LocalRoutes (
	Interface VARCHAR NOT NULL, IP VARCHAR NOT NULL, Gateway VARCHAR NOT NULL, 
	CheckURL VARCHAR, CheckInterval INT NOT NULL DEFAULT 0,
	Weight INT CHECK (Weight >= 0) NOT NULL DEFAULT 1,
	Status VARCHAR CHECK (UPPER(Status) IN ('ONLINE','SHUNNED')) NOT NULL DEFAULT 'ONLINE',
	PRIMARY KEY (IP, Interface, Gateway) );
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func WriteLocalRoutes(db *sql.DB, lr *LocalRoutes) error {
	if lr == nil {
		return errors.New("Empty LR values are given")
	}
	stmt, err := db.Prepare("REPLACE INTO LocalRoutes (Interface, IP, Gateway, CheckURL, CheckInterval, Weight) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(lr.Interface, lr.IP, lr.Gateway, lr.CheckURL, lr.CheckInterval, lr.Weight)
	return err
}

func ReadLocalRoutes(db *sql.DB) ([]LocalRoutes, error) {
	rows, err := db.Query("SELECT Interface, IP, Gateway, CheckURL, CheckInterval, Weight, Status FROM LocalRoutes")
	if err != nil {
		return nil, err
	}

	var lrs []LocalRoutes
	for rows.Next() {
		var lr LocalRoutes
		rows.Scan(&lr.Interface, &lr.IP, &lr.Gateway, &lr.CheckURL, &lr.CheckInterval, &lr.Weight, &lr.Status)
		lrs = append(lrs, lr)
	}
	return lrs, nil
}
