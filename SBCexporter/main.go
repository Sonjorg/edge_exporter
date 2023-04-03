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
	go func() {
	systemResourceCollector()
	wg.Done()

	}()
	go func() {
		hardwareCollector()
	wg.Done()
	}()
	go func() {
		routingEntryCollector()
		wg.Done()
	}()
	go func() {

	callStatsCollector()
	wg.Done()

	}()
	//wg.Wait()


	//wg.Wait()

	//foo := newFooCollector()
	//prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9111", nil))

}
