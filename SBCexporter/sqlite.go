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

func createTable(db *sql.DB) error {
	createAuthTableSQL := `CREATE TABLE IF NOT EXISTS authentication (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"phpsessid" TEXT,
		"time" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create table...")
	statement, err := db.Prepare(createAuthTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	statement.Exec() // Execute SQL Statements
	log.Println("table created")
	return nil
}

func dropTable(db *sql.DB) error{
	dropAuthTableSQL := `DROP TABLE [IF EXISTS] authentication`
	statement, err := db.Prepare(dropAuthTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}// Execute SQL Statements
	fmt.Println("table dropped")
	return nil
}


// We are passing db reference connection from main to our method with other parameters
func insertAuth(db *sql.DB, ipaddress string, phpsessid string, time string) error{
	log.Println("Inserting record ...")
	insertAuthSQL := `INSERT OR REPLACE INTO authentication(ipaddress, phpsessid, time) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertAuthSQL) // Prepare statement.
                                                   // This is good to avoid SQL injections
	if err != nil {
		return err
	}
	_, err = statement.Exec(ipaddress, phpsessid, time)
	if err != nil {
		return err
	}
	return nil
}


func displayOneVal(db *sql.DB, ip string) []*Cookie{
	row, err := db.Query("SELECT * FROM authentication WHERE ipaddress = ? ") //WHERE Start = '" & $start & "'
	//stmt, err := db.Prepare("SELECT * FROM authentication WHERE ipaddress = ?")
	row.Scan(ip)
	if err != nil {
		//log.Fatal(err)
	}
	defer row.Close()
	var c []*Cookie

	for row.Next() {
		p := &Cookie{}
		if err := row.Scan(&p.Id, &p.Ipaddress, &p.Phpsessid, &p.Time); err != nil{
			 fmt.Println(err)
		}
		c = append(c, p)
}
return c

}
func displayAuth(db *sql.DB, ip string) []*Cookie{
	row, err := db.Query("SELECT * FROM authentication")
	//row.Scan(ip)
	if err != nil {
		//log.Fatal(err)
	}
	defer row.Close()

    //checkErr(err)
	/*for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var ipaddress string
		var phpsessid string
		var time string
		row.Scan(&id, &ipaddress, &phpsessid, &time)
		//log.Println("Student: ", code, " ", name, " ", program)
	}*/
	//return row.Columns()ipaddress, phpsessid, time
	//rows.Scan(&id, &name)
	var c []*Cookie
	for row.Next() {
			p := &Cookie{}
			if err := row.Scan(&p.Id, &p.Ipaddress, &p.Phpsessid, &p.Time); err != nil{
				 fmt.Println(err)
			}
			c = append(c, p)
	}
	return c
}
//func displayAuth(db *sql.DB, ip string) *Cookie {

	/*row, err := db.Query("SELECT ipaddress AS ip FROM Students")
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
	}
	//return row.Columns()ipaddress, phpsessid, time
	var c []*Cookie
	for row.Next() {
			p := &Cookie{}
			if err := row.Scan(&p.Id, &p.Ipaddress, &p.Phpsessid, &p.Time); err != nil{
				 fmt.Println(err)
			}
			c = append(c, p)
	}
	return c
	*/
