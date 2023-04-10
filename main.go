package main

import (
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/collector"
	//"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"database/sql"
	"fmt"
	"os"

	//_ "github.com/mattn/go-sqlite3"

)

func main() {

	//Creating database

	var sqliteDatabase *sql.DB
	_, err := os.Stat("sqlite-database.db")
	if err != nil {
		fmt.Println("Creating sqlite-database.db...")
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			fmt.Println(err.Error())
		}
		file.Close()
	}
	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		fmt.Println(err)
	}

	// Creating tables
	err = database.CreateTable(sqliteDatabase)
	if err != nil {
		fmt.Println(err)
	}
	err = database.CreateRoutingSqlite(sqliteDatabase)
	if err != nil {
		fmt.Println(err)
	}

	//starting collectors
	collector.SystemResourceCollector()
	collector.HardwareCollector()
	collector.RoutingEntryCollector()
	collector.CallStatsCollector()

	//Serving metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":1231", nil))

}
