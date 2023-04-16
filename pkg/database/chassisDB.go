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

	//log.Println("Create table...")
	statement, err := db.Prepare(createAuthTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	statement.Exec() // Execute SQL Statements
	//log.Println("table created")
	return nil
}

func InsertChassis(db *sql.DB, ipaddress string, chassisType string, serialNumber string) error{
	//log.Println("Inserting session Chassis data ...")
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
	//row.Scan(ip)
	if err != nil {
		log.Print(err)
		return "","", err
	}
	defer row.Close()

	//var c *Chassis
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
            // a real error happened! you should change your function return
            // to "(bool, error)" and return "false, err" here
            log.Print(err)
        }

        return false
    }

    return true
}