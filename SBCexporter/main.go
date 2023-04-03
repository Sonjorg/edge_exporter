package main

import (

	//"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"

)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	go systemResourceCollector()
	wg.Done()

	go routingEntryCollector()
	wg.Done()

	go callStatsCollector()
	wg.Done()

	go hardwareCollector()
	wg.Done()
	//wg.Wait()

	//foo := newFooCollector()
	//prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9111", nil))

}
