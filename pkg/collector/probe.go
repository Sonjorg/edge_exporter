/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
 */
package collector

import (
	"edge_exporter/pkg/config"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AllCollectors struct{
	metrics []prometheus.Metric
}

func (m *AllCollectors) Probe() {
	cfg, err := config.GetConfig(&config.HostCompose{})
	if err != nil {
		println("probe ", err)
	}

	metrics, success := SystemCollector(cfg)
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	if success {
		if (!config.Excluded("linecard")) {
			metrics = LinecardCollector2(cfg)
			for i := range metrics {
				m.metrics= append(m.metrics, metrics[i])
			}
		}
		if (!config.Excluded("routingentry")) {
			metrics = RoutingEntryCollector(cfg)
			for i := range metrics {
				m.metrics= append(m.metrics, metrics[i])
			}
		}
		if (!config.Excluded("ethernetport")) {
			metrics = EthernetPortCollector(cfg)
			for i := range metrics {
				m.metrics= append(m.metrics, metrics[i])
			}
		}
		if (!config.Excluded("diskpartition")) {

			metrics = DiskPartitionCollector(cfg)
			for i := range metrics {
				m.metrics= append(m.metrics, metrics[i])
			}
		}
		if (!config.Excluded("systemcallstats")) {
			metrics = CallStatsCollector(cfg)
			for i := range metrics {
				m.metrics= append(m.metrics, metrics[i])
			}
		}
	}
}

//Collect implements required collect function for all prometheus collectors
func (collector *AllCollectors) Collect(c chan<- prometheus.Metric) {
	for _, m := range collector.metrics {
		c <- m
	}
}

func (collector *AllCollectors) Describe(ch chan<- *prometheus.Desc) {
//The required Describe interface is empty in this project as it is redundant.
//It was decided to change the overall structure of the code to match a similar approach as the fortigate exporter
//in order to solve several issues regarding the previous design.
//Metrics descriptions are now instead defined directly in each collector.
}

func ProbeHandler(w http.ResponseWriter, r *http.Request) {

	registry := prometheus.NewRegistry()
	pc := &AllCollectors{}
	registry.MustRegister(pc)
	pc.Probe()

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}