package main

import (
	"edge_exporter/pkg/collector"
	//"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"

)

func main() {
	collector.SystemResourceCollector()
	collector.HardwareCollector()
	collector.RoutingEntryCollector()
	collector.CallStatsCollector()






	//foo := newFooCollector()
	//prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9100", nil))

}
