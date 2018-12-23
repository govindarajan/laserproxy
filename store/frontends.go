package store

import (
	"database/sql"
	"errors"
	"net"
)

type ProxyType string

const (
	PrTypeForward ProxyType = "FORWARD"
	PrTypeReverse ProxyType = "REVERSE"
	PrTypeBoth    ProxyType = "BOTH"
)

type RouteType string

const (
	WEIGHT RouteType = "WEIGHT"
	BEST   RouteType = "BEST"
)

type Frontends struct {
	Id         int
	ListenAddr net.IP
	Port       int
	Balance    RouteType
	Type       ProxyType
}

func InitFrontends(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS Frontends (Id INT NOT NULL DEFAULT 0, 
		ListenAddr VARCHAR NOT NULL , Port INT CHECK (Port >= 0) NOT NULL DEFAULT 8080,  
		Balance VARCHAR CHECK (Balance in ('WEIGHT','BEST')) NOT NULL DEFAULT 'BEST',
		Type VARCHAR CHECK (UPPER(Type) IN ('FORWARD','REVERSE','BOTH')) NOT NULL DEFAULT 'BOTH', 
		PRIMARY KEY (Id) );
		`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func WriteFrontends(db *sql.DB, fe *Frontends) error {
	if fe == nil {
		return errors.New("Empty Proxy values are given")
	}
	stmt, err := db.Prepare("REPLACE INTO HTTPSourceProxies (Id, ListenAddr, Port, Balance, Type) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(fe.Id, fe.ListenAddr, fe.Port, fe.Balance, fe.Type)
	return err
}
