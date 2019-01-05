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

type LocalRoute struct {
	Interface     string
	IP            net.IP
	Gateway       net.IP
	CheckURL      *string
	CheckInterval int
	Weight        int
	Status        LRStatus
}

func InitLocalRoute(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS LocalRoute (
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

func WriteLocalRoute(db *sql.DB, lr *LocalRoute) error {
	if lr == nil {
		return errors.New("Empty LR values are given")
	}
	stmt, err := db.Prepare("REPLACE INTO LocalRoute (Interface, IP, Gateway, CheckURL, CheckInterval, Weight) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(lr.Interface, lr.IP, lr.Gateway, lr.CheckURL, lr.CheckInterval, lr.Weight)
	return err
}

func DeleteLocalRoute(db *sql.DB, lr *LocalRoute) error {
	if lr == nil {
		return errors.New("Empty LR values are given")
	}
	stmt, err := db.Prepare("DELETE FROM LocalRoute WHERE Interface = ? AND IP = ? AND Gateway = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(lr.Interface, lr.IP, lr.Gateway)
	return err
}

func ReadLocalRoutes(db *sql.DB) ([]LocalRoute, error) {
	rows, err := db.Query("SELECT Interface, IP, Gateway, CheckURL, CheckInterval, Weight, Status FROM LocalRoute")
	if err != nil {
		return nil, err
	}

	var lrs []LocalRoute
	for rows.Next() {
		var lr LocalRoute
		rows.Scan(&lr.Interface, &lr.IP, &lr.Gateway, &lr.CheckURL, &lr.CheckInterval, &lr.Weight, &lr.Status)
		lrs = append(lrs, lr)
	}
	return lrs, nil
}
