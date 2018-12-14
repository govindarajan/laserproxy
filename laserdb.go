package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// liveRequests struct, represent liveRequests fields
type liveRequests struct {
	id         int
	method     string
	requestURI string
	sourceID   int
	targetID   int
}

// HTTPTargets struct, represent HTTPTargets fields
type HTTPTargets struct {
	id          int
	host        string
	ip          string
	port        int
	weight      int
	status      int
	maxRequests int
}

// requestCounter struct, represent requestCounter fields
type requestCounter struct {
	id         int
	url        string
	method     string
	count      int
	statuscode int
}

// liveRoutes struct, represent liveRoutes fields
type liveRoutes struct {
	id   int
	host string
	ip   string
	port int
}

// HTTPSourceRoutes struct, represent HTTPSourceRoutes fields
type HTTPSourceRoutes struct {
	id              int
	Interface       string
	Type            string
	healthcheck     string
	internalGateway string
	externalGateway string
	weight          int
	active          int
	maxRequests     int
}

func main() {

	// Remove the laserproxy database file if exists. we can comment the
	// line if we dont need to delete
	os.Remove("./laserproxy.db")

	// Initiate database connection
	db, err := sql.Open("sqlite3", "./laserproxy.db")

	// Check if database connection was Initiated successfully
	if err != nil {
		// Print error and exit if there was problem opening connection.
		log.Fatal(err)
	}
	// close database connection before exiting program.
	defer db.Close()

	{ // Create table Block.SQL statement to create all the tables that we need for lazerProxy, with no records in it.
		sqlStmt := `
		CREATE TABLE IF NOT EXISTS liveRequests (id INTEGER PRIMARY KEY AUTOINCREMENT, method TEXT, requestURI TEXT, sourceID INTEGER, targetID INTGER);DELETE FROM liveRequests;
		CREATE TABLE IF NOT EXISTS HTTPTargets (id INTEGER PRIMARY KEY AUTOINCREMENT, host TEXT, ip TEXT, port INTEGER, weight INTEGER, status TEXT, maxRequests INTEGER);DELETE FROM HTTPTargets;
		CREATE TABLE IF NOT EXISTS requestCounter (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, method TEXT, count INTEGER, statusCode INTEGER);DELETE FROM requestCounter;
		CREATE TABLE IF NOT EXISTS liveRoutes (id INTEGER PRIMARY KEY AUTOINCREMENT, host TEXT, ip TEXT, port INTEGER);DELETE FROM liveRoutes;
		CREATE TABLE IF NOT EXISTS HTTPSourceRoutes (id INTEGER PRIMARY KEY AUTOINCREMENT, Interface TEXT, Type TEXT, healthCheck INTEGER, internalGateway TEXT, externalGateway TEXT, weight INTEGER, status TEXT, MaxRequests INTEGER);DELETE FROM HTTPSourceRoutes;
		`
		// Execute the SQL statement
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("Could not create tables %q: %s\n", err, sqlStmt)
			return
		}
	}

	{ // Create records. Begin transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		// Prepare prepared statement that can be reused.
		stmt, err := tx.Prepare("INSERT INTO liveRequests(method, requestURI, sourceID, targetID) VALUES(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		// close statement before exiting program.
		defer stmt.Close()

		// Create empty slice of all struct pointers.
		liveRequests := []*liveRequests{}
		liveRoutes := []*liveRoutes{}
		HTTPSourceRoutes := []*HTTPSourceRoutes{}
		HTTPTargets := []*HTTPTargets{}
		requestCounter := []*requestCounter{}

		// Insert records, Execute statements for each tasks
		for i := range tasks {
			_, err = stmt.Exec(liveRequests[i].method, liveRequests[i].requestURI, liveRequests[i].sourceID, liveRequests[i].targetID)
			if err != nil {
				log.Fatal(err)
			}
		}
		tx.Commit() // Commit the transaction.
	}

	{ // Read records Block, Start reading records
		stmt, err := db.Prepare("SELECT id, method, requestURI from liveRequests where targetID = ?")
		if err != nil {
			log.Fatal(err)
			log.Panic
		}
		defer stmt.Close()

		rows, err := stmt.Query(0)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var method string
			var requestURI string
			err = rows.Scan(&id, &method, &requestURI)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, method, requestURI)
		}
		// To just check if any error was occured during iteration.
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	{ // Update record(s)
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("UPDATE liveRequests SET method = ? where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec("OPTIONS", 2)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	}

	//For dont need since deletion mostly will do from admin tool.
	{ // Delete records block
		// Delete record(s)s
		// _, err = db.Exec("DELETE from task")
		// if err != nil {
		//	log.Fatal(err)
		// }
	}

}
