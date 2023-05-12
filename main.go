package main

import (
	"edge_exporter/pkg/collector"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/utils"
	thishttp "edge_exporter/pkg/http"
	"log"
	"net/http"
)

func main() {
	//Creating database and tables
	database.InitializeDB()

	hosts := config.GetAllHosts()
	for i := range hosts {
		//Fetching sessioncookies and inserting them into the database
		phpsessid, err := thishttp.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err!= nil {
			log.Print(err)
			continue
		}
		//Fetching SBC type and serialnumbers, and inserting them in database
		_, _, err = utils.GetChassisLabels(hosts[i].Ip, phpsessid)
		if err!= nil {
			log.Print(err)
			continue
		}
	}

	http.HandleFunc("/metrics", collector.ProbeHandler)

	log.Fatal(http.ListenAndServe(":5123", nil))

	log.Println("Edge exporter running, listening on 5123")
	select {}
}
