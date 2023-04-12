package utils

//system status exporter
//rest/system/historicalstatistics/1

import (
	"encoding/xml"
	"fmt"
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
	SerialNumber      string `xml:"SerialNumber"`    // Average percent usage of the CPU.
}


//Collect implements required collect function for all promehteus collectors
func GetChassisLabels(ipaddress string, phpsessid string) (chassisType string, serialNumber string, err error){
	//hosts := config.GetAllHosts()//retrieving targets for this exporter

	//fmt.Println(hosts)
	var sqliteDatabase *sql.DB
	var labels *database.Chassis
	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		fmt.Println(err)
		return "","", err
	} // Open the created SQLite File
	// Defer Closing the database
	defer sqliteDatabase.Close()
	if (database.RowExists(sqliteDatabase, ipaddress)) {
		labels, err = database.GetChassis(sqliteDatabase, ipaddress)
	} else {

			dataStr := "https://"+ipaddress+"/rest/chassis"
			_, data,err := http.GetAPIData(dataStr, phpsessid)
			if err != nil {
					fmt.Println("Error collecting from : ", err)

			}
			b := []byte(data) //Converting string of data to bytestream
			ssbc := &ChassisData{}
			xml.Unmarshal(b, &ssbc) //Converting XML data to variables
			//fmt.Println("Successful API call data: ",ssbc.SystemData,"\n")

			chassisType := ssbc.Chassis.Rt_Chassis_Type
			serialNumber := ssbc.Chassis.SerialNumber
			err = database.InsertChassis(sqliteDatabase, ipaddress, chassisType, serialNumber)
				if err != nil {
					fmt.Println("insert chassis error", err)
				}
			return chassisType, serialNumber, err
	}
return labels.ChassisType, labels.SerialNumber, err
}
