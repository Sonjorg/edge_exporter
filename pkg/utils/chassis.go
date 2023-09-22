package utils

//system status exporter
//rest/system/historicalstatistics/1

/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
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
	Rt_Chassis_Type        string `xml:"rt_Chassis_Type"`
	SerialNumber           string `xml:"SerialNumber"`
	CoreSwitch_Temperature int    `xml:"rt_Chassis_CoreSwitch_Temp"`
}


func GetChassisLabelsHTTP(ipaddress string, phpsessid string) (chassisType string, serialNumber string, temperature int, err error){

				dataStr := "https://"+ipaddress+"/rest/chassis"
				_, data,err := http.GetAPIData(dataStr, phpsessid)
				if err != nil {

					return "Error fetching chassisinfo","Error fetching chassisinfo", 0, err
				}
				ssbc := &ChassisData{}
				err = xml.Unmarshal(data, &ssbc) //Converting XML data to variables
				if err != nil {
					return "Error fetching chassisinfo","Error fetching chassisinfo",0, err
				}

				chassisType = ssbc.Chassis.Rt_Chassis_Type
				serialNumber = ssbc.Chassis.SerialNumber
				temperature = ssbc.Chassis.CoreSwitch_Temperature
				//DB
				var sqliteDatabase *sql.DB
				sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
				if err != nil {
					log.Print(err)
				} // Open the created SQLite File
				// Defer Closing the database
				defer sqliteDatabase.Close()
				if (!database.RowExists(sqliteDatabase, ipaddress)) {
					err = database.InsertChassis(sqliteDatabase, ipaddress, chassisType, serialNumber)
						if err != nil {
							log.Print("insert chassis error", err)
					}
				}

return chassisType, serialNumber, temperature, err

}
func GetChassisLabelsDB(ipaddress string) (chassisType string, serialNumber string, err error){

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
					return "Error fetching chassisinfo","Error fetching chassisinfo",err
			}
		}
	
return chassisType, serialNumber, err
}

