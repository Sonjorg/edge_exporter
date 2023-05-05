package collector

import (
	//"fmt"
	//"edge_exporter/pkg/collector"
	"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
)

type AllCollectors struct{
	metrics []prometheus.Metric
}

func (m *AllCollectors) Probe() {
	metrics := LinecardCollector2()
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	metrics = SystemCollector()
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
	metrics = RoutingEntryCollector()
	for i := range metrics {
		m.metrics= append(m.metrics, metrics[i])
	}
}


func (collector *AllCollectors) Collect(c chan<- prometheus.Metric) {
	for _, m := range collector.metrics {
		c <- m
	}
}
func (collector *AllCollectors) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	//ch <- collector.Rt_CardType
	//ch <- collector.Rt_Location
	
	//ch <- collector.Error_ip
}
