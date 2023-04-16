package main

import (
	"edge_exporter/pkg/collector"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	myhttp "edge_exporter/pkg/http"
	//"fmt"

	//"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	//_ "github.com/mattn/go-sqlite3"
)

func main() {


	//Creating database and tables
	database.InitializeDB()
//Fetching sessioncookies
	hosts := config.GetAllHosts()
	for i := range hosts {
		_, err := myhttp.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		fmt.Println("myhttp, err", err)
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
