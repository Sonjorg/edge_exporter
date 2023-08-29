/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
 */
package database

import (
	"database/sql"
	"fmt"
	"log"
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//_ "github.com/mattn/go-sqlite3"
	//"fmt"
)

type RoutingT struct {
	Id        int
	Ipaddress string
	Time      string
	RoutingTable string
	RoutingEntry string
	RtDesc string
}

func CreateRoutingSqlite(db * sql.DB) error{
	createRoutingTables := `CREATE TABLE IF NOT EXISTS routingtables (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"time" TEXT,
		"routingtable" TEXT,
		"routingentries" TEXT,
		"rtdesc" TEXT
		);`
	statement, err := db.Prepare(createRoutingTables)
	if err != nil {
		log.Print(err)
		return err
	}

	statement.Exec()
	return nil
}



func StoreRoutingEntries(db *sql.DB, ipaddress string, time string, routingTable string, routingEntries []string, rtdesc string) error{

	for i := range routingEntries {
		insertSQL1 := `INSERT INTO routingtables(ipaddress, time, routingtable, routingentries, rtdesc) VALUES (?, ?, ?, ?, ?)`

		statement, err := db.Prepare(insertSQL1)
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = statement.Exec(ipaddress, time, routingTable, routingEntries[i], rtdesc)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func GetRoutingData(db *sql.DB,ipaddress string) (map[string][]string,[]string,string, string, error) {
		row, err := db.Query("SELECT * FROM routingtables")
		if err != nil {
			log.Print(err)
			return nil, nil,"","", err
		}
		defer row.Close()

		var time string
		var rtdesc string
		var routingEntries = make(map[string][]string)
		var tables []string

		for row.Next() {
			r := &RoutingT{}
				err := row.Scan(&r.Id, &r.Ipaddress,&r.Time,&r.RoutingTable, &r.RoutingEntry, &r.RtDesc)
				if err != nil{
					log.Print(err)
				}
				fmt.Println("sql ",&r.Id, &r.Ipaddress,&r.Time,&r.RoutingTable, &r.RoutingEntry, &r.RtDesc)
					routingEntries[r.RoutingTable] = append(routingEntries[r.RoutingTable], r.RoutingEntry)
					time = r.Time
					rtdesc = r.RtDesc
		}

		for key, _ := range routingEntries {
			tables = append(tables, key)
		}
		return routingEntries, tables, time, rtdesc, err
}

func DeleteRoutingTables(db *sql.DB, ipaddress string) {
	stmt, _ := db.Prepare("delete from routingtables where ipaddress=?")

    stmt.Exec(ipaddress)
}

func RoutingTablesExists(db * sql.DB,ip string) bool {
	sqlStmt := `SELECT ipaddress FROM routingtables WHERE ipaddress = ?`
    err := db.QueryRow(sqlStmt, ip).Scan(&ip)
    if err != nil {
        if err != sql.ErrNoRows {

            log.Print(err)
        }

        return false
    }

    return true
}