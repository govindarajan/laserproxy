package store

import (
	"database/sql"
	"errors"
)

type TargetBalanceType string

const (
	TBT_BEST TargetBalanceType = "BEST"
	TBT_RR   TargetBalanceType = "RR"
)

type Targets struct {
	Hostname      string
	CheckURL      *string
	CheckInterval int
	Balance       TargetBalanceType
}

func InitTargets(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS Targets (Hostname VARCHAR NOT NULL, CheckURL VARCHAR, 
		CheckInterval INT CHECK (CheckInterval >= 0) NOT NULL DEFAULT 0, Balance VARCHAR CHECK (UPPER(Balance) IN ('BEST', 'RR')),
		PRIMARY KEY (hostname) );
		`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func WriteTargets(db *sql.DB, t *Targets) error {
	if t == nil {
		return errors.New("Empty target values are given")
	}
	stmt, err := db.Prepare("INSERT or REPLACE INTO Targets (Hostname, CheckURL, CheckInterval, Balance) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t.Hostname, t.CheckURL, t.CheckInterval, t.Balance)
	return err
}
