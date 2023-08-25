/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
 */
package main

import (
	"edge_exporter/pkg/collector"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	thishttp "edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"log"
	"net/http"
)

func main() {
	//Creating database and tables
	database.InitializeDB()

	host := config.GetConfig(&config.HostCompose{})
		//Fetching sessioncookies and inserting them into the database
		if (thishttp.SBCIsUp(host.Ip)){
			collector.TestXML(host)

			phpsessid, err := thishttp.APISessionAuth(host.Username, host.Password, host.Ip)
			if err!= nil {
				log.Print(err)
			}
			//Fetching SBC type and serialnumbers, and inserting them in database
			_, _, err = utils.GetChassisLabels(host.Ip, phpsessid)
			if err!= nil {
				log.Print(err)
			}
		}
		
	http.HandleFunc("/metrics", collector.ProbeHandler)
	
	log.Fatal(http.ListenAndServe(":5123", nil))

	log.Println("Edge exporter running, listening on 5123")
	select {}
	
	//collector.TestXML()
}
