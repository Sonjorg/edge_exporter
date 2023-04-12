package database

import (
	"database/sql"
	"log"
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//_ "github.com/mattn/go-sqlite3"
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

func CreateRoutingSqlite(db * sql.DB) error{
	createRoutingTables := `CREATE TABLE IF NOT EXISTS routingtables (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"time" TEXT,
		"routingtable" TEXT,
		"routingentries" TEXT
		);`
fmt.Println("creating routingtables")
	statement, err := db.Prepare(createRoutingTables) // Prepare SQL Statement
	if err != nil {
		return err
	}

	statement.Exec()
	return nil
}

func StoreRoutingEntries(db *sql.DB, ipaddress string, time string, routingTable string, routingEntries []string) error{
	log.Println("Inserting entry ...")
	for i := range routingEntries {
		insertSQL1 := `INSERT OR REPLACE INTO routingtables(ipaddress, time, routingtable, routingentries) VALUES (?, ?, ?, ?)`

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

func RoutingTablesExists(db * sql.DB,ip string) bool {
   // sqlStmt := `SELECT ipaddress FROM routingtables WHERE ipaddress = ?`
	sqlStmt := `SELECT ipaddress FROM routingtables WHERE ipaddress = ?`
    err := db.QueryRow(sqlStmt).Scan(&ip)
    if err != nil {
        if err != sql.ErrNoRows {

            return false
        }

        return false
    }

    return true
}


func GetRoutingData(db *sql.DB,ipaddress string) (map[string][]string,[]string,string, error) {
		row, err := db.Query("SELECT * FROM routingtables WHERE ipaddress = ?", ipaddress)
		//row.Scan(ip)
		if err != nil {
			return nil, nil,"", err
			//fmt.Println(err)
		}

		defer row.Close()

		//var re []RoutingTest
		//var data []*RoutingT
		var time string
		var routingEntries = make(map[string][]string)
		var tables []string

		for row.Next() {
			r := &RoutingT{}
				if err := row.Scan(&r.Id, &r.Ipaddress,&r.Time,&r.RoutingTable, &r.RoutingEntry); err != nil{
					fmt.Println(err)
				}
				if (r.Ipaddress == ipaddress) {
					routingEntries[r.RoutingTable] = append(routingEntries[r.RoutingTable], r.RoutingEntry)
					time = r.Time
				}
		}

		for key, _ := range routingEntries {
			tables = append(tables, key)
		}
		return routingEntries,tables,time,err
	}

/*
func main() {

	var sqliteDatabase *sql.DB

				sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
				if err != nil {
					fmt.Println(err)
				}
	var s []string
	s = append(s, "1")
	s = append(s, "2")
	s = append(s, "3")


	createRoutingSqlite(sqliteDatabase)
	storeRoutingEntries(sqliteDatabase, "ipadresse", "time","5", s)
	if (routingTablesExists(sqliteDatabase, "ipadresse")) {
		g, err := getRoutingEntries(sqliteDatabase,"ipadresse","5")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(g)

	}
}
*/