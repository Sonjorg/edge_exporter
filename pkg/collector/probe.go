/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type AllCollectors struct{
	metrics []prometheus.Metric
}

func (m *AllCollectors) Probe() {
	metrics, successfulHosts := SystemCollector()
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	metrics = LinecardCollector2(successfulHosts)
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	metrics = RoutingEntryCollector(successfulHosts)
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	metrics = EthernetPortCollector(successfulHosts)
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	metrics = DiskPartitionCollector(successfulHosts)
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	metrics = CallStatsCollector(successfulHosts)
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
}

//Collect implements required collect function for all prometheus collectors
func (collector *AllCollectors) Collect(c chan<- prometheus.Metric) {
	for _, m := range collector.metrics {
		c <- m
	}
}

func (collector *AllCollectors) Describe(ch chan<- *prometheus.Desc) {
//The required Describe interface is empty in this project,
//because we decided to change the overall structure of the code to match a similar approach as the fortigate exporter
//in order to solve several issues we had with the previous design.
//Metrics descriptions are now instead defined directly in each collector function
}

func ProbeHandler(w http.ResponseWriter, r *http.Request) {

	registry := prometheus.NewRegistry()
	pc := &AllCollectors{}
	registry.MustRegister(pc)
	pc.Probe()

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}