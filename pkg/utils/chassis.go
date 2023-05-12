package utils

//system status exporter
//rest/system/historicalstatistics/1

import (
	"encoding/xml"
	"log"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"database/sql"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/http"
)

type ChassisData struct {
	XMLname    xml.Name   `xml:"root"`
	Chassis chassis       `xml:"chassis"`
}

type chassis struct {
	Rt_Chassis_Type   string `xml:"rt_Chassis_Type"`
	SerialNumber      string `xml:"SerialNumber"`   
}


func GetChassisLabels(ipaddress string, phpsessid string) (chassisType string, serialNumber string, err error){

	var sqliteDatabase *sql.DB
	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Print(err)
		return "","", err
	} // Open the created SQLite File
	// Defer Closing the database
	defer sqliteDatabase.Close()
	if (database.RowExists(sqliteDatabase, ipaddress)) {
		chassisType, serialNumber, err = database.GetChassis(sqliteDatabase, ipaddress)
		if (chassisType == "" || serialNumber == "" || err != nil) {
			if (phpsessid != "null") {
				dataStr := "https://"+ipaddress+"/rest/chassis"
				_, data,err := http.GetAPIData(dataStr, phpsessid)
				if err != nil {

					return "http error","http error",err
				}
				ssbc := &ChassisData{}
				err = xml.Unmarshal(data, &ssbc) //Converting XML data to variables
				if err != nil {
					return "http error","http error",err
				}

				chassisType := ssbc.Chassis.Rt_Chassis_Type
				serialNumber := ssbc.Chassis.SerialNumber

				err = database.InsertChassis(sqliteDatabase, ipaddress, chassisType, serialNumber)
					if err != nil {
						log.Print("insert chassis error", err)
					}
				return string(chassisType), string(serialNumber), err
			}
		}
	}
return chassisType, serialNumber, err
}
