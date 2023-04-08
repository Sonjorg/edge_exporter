package database

import (
	//"database/sql"
	//"log"
	"time"
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//"fmt"
)
func checkTime(hours int, timeLast time.Time) bool{
	//mins := time.Minute * time.Duration(8)
timeSchedule := time.Hour * time.Duration(hours)
now := time.Now().Format(time.RFC3339)
timeNow, _ := time.Parse(time.RFC3339, now)

if (timeLast.Add(timeSchedule).After(timeNow)) {
	return true
} else {
	 return false}
}