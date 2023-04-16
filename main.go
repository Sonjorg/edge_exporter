package main

import (
	"edge_exporter/pkg/collector"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/utils"
	myhttp "edge_exporter/pkg/http"
	"fmt"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	//Creating database and tables
	database.InitializeDB()

	hosts := config.GetAllHosts()
	for i := range hosts {
		//Fetching sessioncookies and placing them in database
		phpsessid, err := myhttp.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err!= nil {
			log.Print(err)
		}
		//Fetching SBC type and serialnumbers, and placing them in database
		_, _, err = utils.GetChassisLabels(hosts[i].Ip, phpsessid)
		if err!= nil {
			log.Print(err)
		}
	}

	//starting collectors
	collector.SystemResourceCollector()
	collector.DiskPartitionCollector()
	collector.RoutingEntryCollector()
	collector.CallStatsCollector()
	collector.LinecardCollector()
	//collector.EthernetportCollector()

	//Serving metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":5123", nil))

}
