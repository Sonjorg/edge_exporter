package main

import (
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/collector"
	//"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"


	//_ "github.com/mattn/go-sqlite3"

)

func main() {

	//Creating database
	database.InitializeDB()
	//starting collectors
	collector.SystemResourceCollector()
	collector.DiskPartitionCollector()
	collector.RoutingEntryCollector()
	collector.CallStatsCollector()

	//Serving metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":5123", nil))

}
