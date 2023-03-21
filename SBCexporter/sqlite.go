package main

import (
	"database/sql"
	"log"
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"fmt"
)

type Cookie struct {
	Id        int
	Ipaddress string
	Phpsessid string
	Time      string
}

func createTable(db *sql.DB) {
	createAuthTableSQL := `CREATE TABLE authentication (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"phpsessid" TEXT,
		"time" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create table...")
	statement, err := db.Prepare(createAuthTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertAuth(db *sql.DB, ipaddress string, phpsessid string, time string) {
	log.Println("Inserting record ...")
	insertAuthSQL := `INSERT INTO authentication(ipaddress, phpsessid, time) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertAuthSQL) // Prepare statement.
                                                   // This is good to avoid SQL injections
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(ipaddress, phpsessid, time)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayAuth(db *sql.DB) []*Cookie{
	row, err := db.Query("SELECT * FROM authentication")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	/*for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var ipaddress string
		var phpsessid string
		var time string
		row.Scan(&id, &ipaddress, &phpsessid, &time)
		//log.Println("Student: ", code, " ", name, " ", program)
	}*/
	//return row.Columns()ipaddress, phpsessid, time
	var c []*Cookie
	for row.Next() {
			p := &Cookie{}
			if err := row.Scan(p.Id, p.Ipaddress, p.Phpsessid, p.Time); err != nil{
				 fmt.Println(err)
			}
			c = append(c, p)
	}
	return c
}