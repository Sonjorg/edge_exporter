package main

import (
	//"fmt"
	"edge_exporter/pkg/collector"
	"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
)

type AllCollectors struct{
	metrics []prometheus.Metric
}

func (m *AllCollectors) Probe() {
	metrics := collector.LinecardCollector2()
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
}


func (collector *AllCollectors) Collect(c chan<- prometheus.Metric) {
	for _, m := range collector.metrics {
		c <- m
	}
}