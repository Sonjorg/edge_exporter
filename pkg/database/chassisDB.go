/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
package database

import (
	//"encoding/xml"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"database/sql"
	//"edge_exporter/pkg/http"
)

type Chassis struct {
	Id           int
	Ipaddress    string
	ChassisType  string
	SerialNumber string
}

func CreateChassis(db *sql.DB) error {
	createAuthTableSQL := `CREATE TABLE IF NOT EXISTS chassis (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"chassistype" TEXT,
		"serialnumber" TEXT
	  );` // SQL Statement for Create Table

	statement, err := db.Prepare(createAuthTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	statement.Exec() // Execute SQL Statements
	return nil
}

func InsertChassis(db *sql.DB, ipaddress string, chassisType string, serialNumber string) error{
	log.Println("Inserting chassis data ...")
	insertAuthSQL := `INSERT INTO chassis(ipaddress, chassistype, serialnumber) VALUES (?, ?, ?)`

	statement, err := db.Prepare(insertAuthSQL) // Prepare statement.
                                                   // This is good to avoid SQL injections
	if err != nil {
		return err
	}
	_, err = statement.Exec(ipaddress, chassisType, serialNumber)
	if err != nil {
		return err
	}
	return nil
}

func GetChassis(db *sql.DB, ipaddress string) (string, string, error){
	row, err := db.Query("SELECT * FROM chassis WHERE ipaddress = ?", ipaddress)
	if err != nil {
		log.Print(err)
		return "","", err
	}
	defer row.Close()

	for row.Next() {
			p := &Chassis{}
			if err := row.Scan(&p.Id, &p.Ipaddress, &p.ChassisType, &p.SerialNumber); err != nil{
				 fmt.Println("No data from db", err)
			}
			if (p.Ipaddress == ipaddress) {
				return p.ChassisType, p.SerialNumber,err
			}
	}
return "", "", err

}

func chassisExists(db * sql.DB, ip string) bool {
    sqlStmt := `SELECT ipaddress FROM chassis WHERE ipaddress = ?`
    err := db.QueryRow(sqlStmt, ip).Scan(&ip)
    if err != nil {
        if err != sql.ErrNoRows {
            log.Print(err)
        }

        return false
    }

    return true
}