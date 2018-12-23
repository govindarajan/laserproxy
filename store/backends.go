package store

import (
	"database/sql"
	"errors"
)

type ProxyStatus string

const (
	PrStOnline  ProxyStatus = "ONLINE"
	PrStShunned ProxyStatus = "SHUNNED"
)

type Backends struct {
	GroupId       int
	Hostname      string
	Port          int
	CheckURL      *string
	CheckInterval int
	Weight        int
	Status        ProxyStatus
	MaxRequests   int
}

func InitBackends(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS Backends (GroupId INT NOT NULL DEFAULT 0, 
		Hostname VARCHAR NOT NULL , Port INT CHECK (Port >= 0) NOT NULL DEFAULT 8080, CheckURL VARCHAR, 
		CheckInterval INT NOT NULL DEFAULT 0,  Weight INT CHECK (Weight >= 0) NOT NULL DEFAULT 1,
		Status VARCHAR CHECK (UPPER(Status) IN ('ONLINE','SHUNNED')) NOT NULL DEFAULT 'ONLINE',
		MaxRequests INT CHECK (MaxRequests >=0) NOT NULL DEFAULT 100, 
		PRIMARY KEY (GroupId, Hostname, Port) );
		`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func WriteHTTPSrcProxies(db *sql.DB, be *Backends) error {
	if be == nil {
		return errors.New("Empty Proxy values are given")
	}
	stmt, err := db.Prepare("REPLACE INTO HTTPSourceProxies (GroupId, Hostname, Port, CheckURL, CheckInterval, Weight, MaxRequests) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(be.GroupId, be.Hostname, be.Port, be.CheckURL, be.CheckInterval, be.Weight, be.MaxRequests)
	return err
}
