package database

import (
	"database/sql"
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
//log.Print("creating routingtables")
	statement, err := db.Prepare(createRoutingTables) // Prepare SQL Statement
	if err != nil {
		log.Print(err)
		return err
	}

	statement.Exec()
	return nil
}



func StoreRoutingEntries(db *sql.DB, ipaddress string, time string, routingTable string, routingEntries []string) error{

	for i := range routingEntries {
		insertSQL1 := `INSERT INTO routingtables(ipaddress, time, routingtable, routingentries) VALUES (?, ?, ?, ?)`

		statement, err := db.Prepare(insertSQL1)
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = statement.Exec(ipaddress, time, routingTable, routingEntries[i])
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func DeleteRoutingTables(db *sql.DB, ipaddress string) {
	stmt, err := db.Prepare("delete from routingtables where ipaddress=?")

    res, err := stmt.Exec(ipaddress)

    affect, err := res.RowsAffected()
    log.Println(affect,err)
}

func RoutingTablesExists(db * sql.DB,ip string) bool {
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


func GetRoutingData(db *sql.DB,ipaddress string) (map[string][]string,[]string,string, error) {
		row, err := db.Query("SELECT * FROM routingtables")
		if err != nil {
			log.Print(err)
			return nil, nil,"", err
		}
		defer row.Close()

		var time string
		var routingEntries = make(map[string][]string)
		var tables []string

		for row.Next() {
			r := &RoutingT{}
				err := row.Scan(&r.Id, &r.Ipaddress,&r.Time,&r.RoutingTable, &r.RoutingEntry)
				if err != nil{
					log.Print(err)
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
