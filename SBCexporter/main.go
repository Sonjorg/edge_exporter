package main

import (

	//"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	//"sync"

)

func main() {
	//var wg sync.WaitGroup
	go func() {
	systemResourceCollector()
	}()
	go func() {
		go hardwareCollector()
	//wg.Done()
	}()
	go func() {
		goroutingEntryCollector()
	//	wg.Done()
	}()
	go func() {

	callStatsCollector()
	}()
	//wg.Wait()


	//wg.Wait()

	//foo := newFooCollector()
	//prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9111", nil))

}
