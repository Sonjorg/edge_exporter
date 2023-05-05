package main

import (
	"edge_exporter/pkg/collector"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/utils"
	thishttp "edge_exporter/pkg/http"
	//"fmt"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
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
		//Fetching SBC type and serialnumbers, and placing them in database
		_, _, err = utils.GetChassisLabels(hosts[i].Ip, phpsessid)
		if err!= nil {
			log.Print(err)
			continue
		}
	}

	//initializing collectors
	//The collectors only start to run after a scrape is made to this host
	//collector.SystemResourceCollector()
	//collector.DiskPartitionCollector()
	//collector.RoutingEntryCollector()
	//collector.CallStatsCollector()
	//collector.LinecardCollector()
	//collector.EthernetportCollector()
	registry := prometheus.NewRegistry()
	c := &collector.AllCollectors{}

	registry.MustRegister(c)
	go func() {
		c.Probe()
	}()
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	//h.ServeHTTP()
	//Serving metrics
	http.Handle("/metrics", h)
	//http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9103", nil))
	
	log.Printf("Edge exporter running, listening on :9103")
	select {}
}
