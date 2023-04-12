package database

import (
	"database/sql"
	//"log"
	"time"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"fmt"
	"os"
)




func WithinTime(hours float64, previoust time.Time) bool{
	//mins := time.Minute * time.Duration(8)
	var duration time.Duration = time.Duration(hours)

    // in hours
    fmt.Println(duration.Hours())
timeSchedule :=  duration
now := time.Now().Format(time.RFC3339)
timeNowParsed, err := time.Parse(time.RFC3339, now)
if err != nil {
	fmt.Println(err)
	return false
}
if err != nil {
	fmt.Println(err)
	return false
}
return previoust.Add(timeSchedule).After(timeNowParsed)

}

func InitializeDB() {
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
	err = CreateTable(sqliteDatabase)
	if err != nil {
		fmt.Println(err)
	}
	err = CreateRoutingSqlite(sqliteDatabase)
	if err != nil {
		fmt.Println(err)
	}
	err = CreateChassis(sqliteDatabase)
	if err != nil {
		fmt.Println("Chassis DB error",err)
	}

}