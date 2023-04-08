package database

import (
	"database/sql"
	"log"
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"fmt"
)

type RoutingT struct {
	Id        int
	Ipaddress string
	Time      string
	RoutingTable string
	RoutingEntry string
	 //map consisting of routingtables and their routingentries
}
/*
type RoutingE struct {
	Time      string
	RoutingTable string
	RoutingEntry string
}*/
/*
type RoutingTmp struct {
	Id        int
	Ipaddress string
	Time      string
	RoutingTablesnEntries map[string][]string
	//RoutingEntries []string
}*/
func createRoutingSqlite(db * sql.DB) error{
	createRoutingTables := `CREATE TABLE IF NOT EXISTS routingtables (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"time" TEXT,
		"routingtable" TEXT,
		"routingentries" TEXT
		);`

	statement, err := db.Prepare(createRoutingTables) // Prepare SQL Statement
	if err != nil {
		return err
	}

	statement.Exec()
	return nil
}

func storeRoutingTables(db *sql.DB, ipaddress string, time string, routingTable string, routingEntries []string) error{
	log.Println("Inserting record ...")
	for i := range routingEntries {
		insertSQL1 := `INSERT INTO routingtables(ipaddress, time, routingtable, routingentries) VALUES (?, ?, ?, ?)`

		statement, err := db.Prepare(insertSQL1) // Prepare statement.
													// This is good to avoid SQL injections
		if err != nil {
			fmt.Println(err)

			return err

		}
		_, err = statement.Exec(ipaddress, time, routingTable, routingEntries[i])
		if err != nil {
			fmt.Println(err)
			return err
		}
}
	return nil
}

func storeRoutingEntries(db *sql.DB, ipaddress string, time string, routingTable string, routingEntries []string) error{
	log.Println("Inserting record ...")
	for i := range routingEntries {
		insertSQL1 := `INSERT INTO routingtables(ipaddress, time, routingtable, routingentries) VALUES (?, ?, ?, ?)`

		statement, err := db.Prepare(insertSQL1) // Prepare statement.
													// This is good to avoid SQL injections
		if err != nil {
			fmt.Println(err)

			return err

		}
		_, err = statement.Exec(ipaddress, time, routingTable, routingEntries[i])
		if err != nil {
			fmt.Println(err)
			return err
		}
}
	return nil
}

func routingTablesExists(db * sql.DB, ip string) bool {
    sqlStmt := `SELECT ipaddress FROM routingtables WHERE ipaddress = ?`
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



func getRoutingEntries(db *sql.DB,ipaddress string,routingTable string) ([]string, error) {

	//if (routingTablesExists(db,ipaddress)) {
		//row, err := db.Query("SELECT * FROM routingtables")
		row, err := db.Query("SELECT * FROM routingtables")
		//row.Scan(ip)
		if err != nil {
			return nil, err
			//fmt.Println(err)
		}
		defer row.Close()
		var re []string
		//var data []*RoutingT
		for row.Next() {
			r := &RoutingT{}
				if err := row.Scan(&r.Id, &r.Ipaddress,&r.Time,&r.RoutingTable, &r.RoutingEntry); err != nil{
					fmt.Println(err)
				}
				if (r.Ipaddress == ipaddress) {
					//data = append(data, r)
					if (r.RoutingTable == routingTable) {
						re = append(re, r.RoutingEntry)
					}
				}
		}

		return re ,err
	//} else {
	//	return nil, nil
	//}
}


