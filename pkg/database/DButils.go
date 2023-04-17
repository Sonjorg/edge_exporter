package database

import (
	"database/sql"
	"log"
	"time"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"fmt"
	"os"
)


func Expired(hours float64, previoustime time.Time) bool{
	//mins := time.Minute * time.Duration(8)
	var timeSchedule time.Duration = time.Duration(hours)
	duration := timeSchedule*time.Minute
	//log.Print("duration", duration)
	now := time.Now().Format(time.RFC3339)
	timeNowParsed, err := time.Parse(time.RFC3339, now)
	if err != nil {
		log.Print(err)
		return false
	}
	if err != nil {
		log.Print(err)
		return false
	}
	//previoustime.
	//If previous time + 24 hours is not after (before) now: not expired
	fmt.Print(previoustime.Add(duration).Before(timeNowParsed))
	return previoustime.Add(duration).Before(timeNowParsed) //after, contioue using this
}

func InitializeDB() {
	var sqliteDatabase *sql.DB

	_, err := os.Stat("sqlite-database.db")
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
		log.Print("Create routing entry table")

	}
	err = CreateChassis(sqliteDatabase)
	if err != nil {
		log.Print("Chassis DB error",err)
		log.Print("Create chassis table")
	}

}