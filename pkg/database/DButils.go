/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
package database

import (
	"database/sql"
	"log"
	//"time"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//"fmt"
	"os"
)



func InitializeDB() {

	err := os.Remove("sqlite-database.db")
    if err != nil {
        log.Println(err)
    }

	var sqliteDatabase *sql.DB

	_, err = os.Stat("sqlite-database.db")
	if err != nil {
		log.Print("Creating sqlite-database.db...") //"./pkg/database/sqlite-database.db"
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			log.Print(err.Error())
		}
		file.Close()
	}
	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Print(err)
	}

	// Creating tables
	err = CreateAuthTable(sqliteDatabase)
	if err != nil {
		log.Print(err)
	}
	log.Print("Create authentication table")
	err = CreateRoutingSqlite(sqliteDatabase)
	if err != nil {
		log.Print(err)
	}
	log.Print("Create routing entry table")
	err = CreateChassis(sqliteDatabase)
	if err != nil {
		log.Print("Chassis DB error",err)
	}
	log.Print("Create chassis table")

}