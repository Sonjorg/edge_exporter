/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
 */
package collector

import (
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"encoding/xml"
	"log"
	"github.com/prometheus/client_golang/prometheus"
)

type cSBCdata struct {
	XMLname       xml.Name      `xml:"root"`
	Status        cStatus       `xml:"status"`
	CallStatsData callStatsData `xml:"systemcallstats"`
}
type cStatus struct {
	HTTPcode string `xml:"http_code"`
}
type callStatsData struct {
	Href                       string `xml:"href,attr"`
	Rt_NumCallAttempts         int    `xml:"rt_NumCallAttempts"`  // Average percent usage of the CPU.
	Rt_NumCallSucceeded        int    `xml:"rt_NumCallSucceeded"` // Average percent usage of system memory. int
	Rt_NumCallFailed           int    `xml:"rt_NumCallFailed"`
	Rt_NumCallCurrentlyUp      int    `xml:"rt_NumCallCurrentlyUp"`      //Number of currently connected calls system wide. int
	Rt_NumCallAbandonedNoTrunk int    `xml:"rt_NumCallAbandonedNoTrunk"` //Number of rejected calls due to no channel available system wide since system came up.
	Rt_NumCallUnAnswered       int    `xml:"rt_NumCallUnAnswered"`
}

func CallStatsCollector(host *config.HostCompose) (m []prometheus.Metric) {
	
	var (
		Rt_NumCallAttempts = prometheus.NewDesc("edge_callstats_NumCallAttempts",
			"Total number of call attempts system wide since system came up.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_NumCallSucceeded = prometheus.NewDesc("edge_callstats_NumCallSucceeded",
			"Total number of successful calls system wide since system came up.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_NumCallFailed = prometheus.NewDesc("edge_callstats_NumCallFailed",
			"Total number of failed calls system wide since system came up.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_NumCallCurrentlyUp = prometheus.NewDesc("edge_callstats_NumCallCurrentlyUp",
			"Number of currently connected calls system wide.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_NumCallAbandonedNoTrunk = prometheus.NewDesc("edge_callstats_NumCallAbandonedNoTrunk",
			"Number of rejected calls due to no channel available system wide since system came up.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_NumCallUnAnswered = prometheus.NewDesc("edge_callstats_NumCallUnAnswered",
			"Number of unanswered calls system wide since system came up.",
			[]string{"hostip", "hostname"}, nil,
		)
	)

		phpsessid, err := http.APISessionAuth(host.Username, host.Password, host.Ip)
		if err != nil {
			log.Print("Error retrieving session cookie = ", err, "\n")
			return 
		}

		dataStr := "https://" + host.Ip + "/rest/systemcallstats"
		_, data, err := http.GetAPIData(dataStr, phpsessid)
		if err != nil {
			log.Print("Error collecting systemcall data = ", err, "\n")
			return
		}

		ssbc := &cSBCdata{}
		err = xml.Unmarshal(data, &ssbc) //Converting XML data to variables
		if err != nil {
			log.Print("XML error callstats", err)
			return
		}

		metricValue1 := float64(ssbc.CallStatsData.Rt_NumCallAttempts)
		metricValue2 := float64(ssbc.CallStatsData.Rt_NumCallSucceeded)
		metricValue3 := float64(ssbc.CallStatsData.Rt_NumCallFailed)
		metricValue4 := float64(ssbc.CallStatsData.Rt_NumCallCurrentlyUp)
		metricValue5 := float64(ssbc.CallStatsData.Rt_NumCallAbandonedNoTrunk)
		metricValue6 := float64(ssbc.CallStatsData.Rt_NumCallUnAnswered)

		m = append(m, prometheus.MustNewConstMetric(Rt_NumCallAttempts, prometheus.GaugeValue, metricValue1, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_NumCallSucceeded, prometheus.GaugeValue, metricValue2, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_NumCallFailed, prometheus.GaugeValue, metricValue3, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_NumCallCurrentlyUp, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_NumCallAbandonedNoTrunk, prometheus.GaugeValue, metricValue5, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_NumCallUnAnswered, prometheus.GaugeValue, metricValue6, host.Ip, host.Hostname))
	
	return m
}
