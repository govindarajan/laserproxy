package store

import (
	"database/sql"
	"errors"
)

type ProxyStatus string

const (
	PrStOnline  ProxyStatus = "ONLINE"
	PrStShunned ProxyStatus = "SHUNNED"
	PrStOffline ProxyStatus = "OFFLINE"
)

type Backend struct {
	GroupId       int
	Host          string // IP:Port format
	CheckURL      *string
	CheckInterval int
	Weight        int
	Status        ProxyStatus
	MaxRequests   int
}

func InitBackend(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS Backend (GroupId INT NOT NULL DEFAULT 0, 
		Host VARCHAR NOT NULL ,  CheckURL VARCHAR, 
		CheckInterval INT NOT NULL DEFAULT 0,  Weight INT CHECK (Weight >= 0) NOT NULL DEFAULT 1,
		Status VARCHAR CHECK (UPPER(Status) IN ('ONLINE','SHUNNED','OFFLINE')) NOT NULL DEFAULT 'ONLINE',
		MaxRequests INT CHECK (MaxRequests >=0) NOT NULL DEFAULT 100, 
		PRIMARY KEY (GroupId, Host) );
		`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func WriteBackend(db *sql.DB, be *Backend) error {
	if be == nil {
		return errors.New("Empty Proxy values are given")
	}
	stmt, err := db.Prepare("REPLACE INTO Backend (GroupId, Host, CheckURL, CheckInterval, Weight, MaxRequests) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(be.GroupId, be.Host, be.CheckURL, be.CheckInterval, be.Weight, be.MaxRequests)
	return err
}

// ReadBackends used to read all the backends from given frontend id
func ReadBackends(db *sql.DB, gID int) ([]Backend, error) {
	rows, err := db.Query("SELECT GroupId, Host, CheckURL, CheckInterval, Weight, MaxRequests FROM Backend WHERE GroupId = ?", gID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Backend
	for rows.Next() {
		var be Backend
		rows.Scan(&be.GroupId, &be.Host, &be.CheckURL, &be.CheckInterval,
			&be.Weight, &be.MaxRequests)
		res = append(res, be)
	}
	return res, nil
}

// ReadAllBackends used to read all the backends from given DB
func ReadAllBackends(db *sql.DB) ([]Backend, error) {
	rows, err := db.Query("SELECT GroupId, Host, CheckURL, CheckInterval, Weight, MaxRequests FROM Backend")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Backend
	for rows.Next() {
		var be Backend
		rows.Scan(&be.GroupId, &be.Host, &be.CheckURL, &be.CheckInterval,
			&be.Weight, &be.MaxRequests)
		res = append(res, be)
	}
	return res, nil
}
