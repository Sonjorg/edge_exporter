package database

import (
	//"database/sql"
	//"log"
	"time"
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//"fmt"
)
func WithinTime(hours int, previousTime string) bool{
	//mins := time.Minute * time.Duration(8)
timeSchedule := time.Minute * time.Duration(hours)
now := time.Now().Format(time.RFC3339)
timeNowParsed, _ := time.Parse(time.RFC3339, now)
pt,err := time.Parse(time.RFC3339, previousTime)
if err != nil {
	return false
}
if (pt.Add(timeSchedule).After(timeNowParsed)) {
	return true
} else {
	 return false}
}
