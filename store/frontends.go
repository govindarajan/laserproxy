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
)

type RouteType string

const (
	WEIGHT RouteType = "WEIGHT"
	BEST   RouteType = "BEST"
)

type Frontend struct {
	Id         int
	ListenAddr net.IP
	Port       int
	Balance    RouteType
	Type       ProxyType
}

func InitFrontend(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS Frontend (Id INT PRIMARY KEY UNIQUE, 
		ListenAddr VARCHAR NOT NULL , Port INT CHECK (Port >= 0) NOT NULL DEFAULT 8080,  
		Balance VARCHAR CHECK (Balance in ('WEIGHT','BEST')) NOT NULL DEFAULT 'BEST',
		Type VARCHAR CHECK (UPPER(Type) IN ('FORWARD','REVERSE','BOTH')) NOT NULL DEFAULT 'BOTH', 
		UNIQUE(ListenAddr, Port));
		`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func WriteFrontend(db *sql.DB, fe *Frontend) error {
	if fe == nil {
		return errors.New("Empty Proxy values are given")
	}
	stmt, err := db.Prepare("INSERT INTO Frontend (Id, ListenAddr, Port, Balance, Type) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(fe.Id, fe.ListenAddr, fe.Port, fe.Balance, fe.Type)
	return err
}

// ReadFrontends used to read all the front ends from given DB
func ReadFrontends(db *sql.DB) ([]Frontend, error) {
	rows, err := db.Query("SELECT Id, ListenAddr, Port, Balance, Type FROM Frontend")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Frontend
	for rows.Next() {
		var fe Frontend
		rows.Scan(&fe.Id, &fe.ListenAddr, &fe.Port, &fe.Balance, &fe.Type)
		// HACK: Fix it by overriding Scan
		fe.ListenAddr = net.ParseIP(string(fe.ListenAddr))
		res = append(res, fe)
	}
	return res, nil
}
