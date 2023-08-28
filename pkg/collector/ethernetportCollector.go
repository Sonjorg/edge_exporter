/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
package collector

import (
	"encoding/xml"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"log"
	"github.com/prometheus/client_golang/prometheus"
)

// /rest/ethernetport
type eSBCdata struct {
	XMLname       xml.Name      `xml:"root"`
	EthernetData  ethernetData  `xml:"ethernetport"`
}
type ethernetData struct {
Href                          string `xml:"href,attr"`
IfName		                  string `xml:"ifName"`
IfAlias                       string `xml:"ifAlias"`
Rt_ifInBroadcastPkts		  int    `xml:"rt_ifInBroadcastPkts"`
Rt_ifInDiscards		          int    `xml:"rt_ifInDiscards"`
Rt_ifInErrors		          int    `xml:"rt_ifInErrors"`
Rt_ifInMulticastPkts		  int    `xml:"rt_ifInMulticastPkts"`
Rt_ifInOverSizedPkts		  int    `xml:"rt_ifInOverSizedPkts"`
Rt_ifInUcastPkts		      int    `xml:"rt_ifInUcastPkts"`
Rt_ifInUndersizedPkts		  int    `xml:"rt_ifInUndersizedPkts"`
Rt_ifOutBroadcastPkts		  int    `xml:"rt_ifOutBroadcastPkts"`//Displays the number of transmitted broadcast packets on this port.
Rt_ifOutDiscards		      int    `xml:"rt_ifOutDiscards"` //Displays the number of discard errors detected on this port.
Rt_ifOutErrors		          int    `xml:"rt_ifOutErrors"` //Displays the number of errors detected on this port.
Rt_ifOutMulticastPkts		  int    `xml:"rt_ifOutMulticastPkts"` //Displays the number of transmitted multicast packets on this port.
Rt_ifOutUcastPkts	    	  int    `xml:"rt_ifOutUcastPkts"` //Displays the number of transmitted unicast packets on this port.
}

func EthernetPortCollector(host *config.HostCompose)(m []prometheus.Metric) {

var (

		Rt_ifInBroadcastPkts = prometheus.NewDesc("edge_ethernet_ifInBroadcastPkts",
			"Displays the number of received broadcast packets on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifInDiscards = prometheus.NewDesc("edge_ethernet_ifInDiscards",
			"Displays the number of discard errors detected on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifInErrors = prometheus.NewDesc("edge_ethernet_ifInErrors",
			"Displays the number of errors detected on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifInMulticastPkts = prometheus.NewDesc("edge_ethernet_ifInMulticastPkts",
			"Displays the number of received multicast packets on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifInOverSizedPkts = prometheus.NewDesc("edge_ethernet_ifInOverSizedPkts",
			"Displays the number of Oversized Packet errors detected on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifInUcastPkts = prometheus.NewDesc("edge_ethernet_ifInUcastPkts",
			"Displays the number of received unicast packets on this port.  ",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifInUndersizedPkts = prometheus.NewDesc("edge_ethernet_ifInUndersizedPkts",
			"Displays the number of Undersized Packet errors detected on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifOutBroadcastPkts = prometheus.NewDesc("edge_ethernet_ifOutBroadcastPkts",
			"Displays the number of transmitted broadcast packets on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifOutDiscards = prometheus.NewDesc("edge_ethernet_ifOutDiscards",
			"Displays the number of discard errors detected on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifOutErrors = prometheus.NewDesc("edge_ethernet_ifOutErrors",
			"Displays the number of errors detected on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifOutMulticastPkts = prometheus.NewDesc("edge_ethernet_ifOutMulticastPkts",
			"Displays the number of transmitted multicast packets on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifOutUcastPkts = prometheus.NewDesc("edge_ethernet_ifOutUcastPkts",
			"Displays the number of transmitted unicast packets on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
)

		phpsessid,err := http.APISessionAuth(host.Username, host.Password, host.Ip)
		if err != nil {
			log.Print("Error session cookie: ", host.Ip, err)
			return
		}
		var ethernetportID []string
			//Every router has these ethernetports regardless of SBC1000 or SBC2000, according to HDO
			ethernetportID = append(ethernetportID, "23")
			ethernetportID = append(ethernetportID, "29")
			ethernetportID = append(ethernetportID, "24")
			for j := range ethernetportID {
					url := "https://"+host.Ip+"/rest/ethernetport/"+ethernetportID[j]
					_, data, err := http.GetAPIData(url, phpsessid)
						if err != nil {
							log.Print(err)
							continue
						}
					eData := &eSBCdata{}
					err = xml.Unmarshal(data, &eData) //Converting XML data to variables
					if err!= nil {
						log.Print("XML error ethernet", err)
						continue
					}

					metricValue3 := float64(eData.EthernetData.Rt_ifInBroadcastPkts)
					metricValue4 := float64(eData.EthernetData.Rt_ifInDiscards)
					metricValue5 := float64(eData.EthernetData.Rt_ifInErrors)
					metricValue8 := float64(eData.EthernetData.Rt_ifInMulticastPkts)
					metricValue10 := float64(eData.EthernetData.Rt_ifInOverSizedPkts)
					metricValue11 := float64(eData.EthernetData.Rt_ifInUcastPkts)
					metricValue12 := float64(eData.EthernetData.Rt_ifInUndersizedPkts)
					metricValue19 := float64(eData.EthernetData.Rt_ifOutBroadcastPkts)
					metricValue21 := float64(eData.EthernetData.Rt_ifOutDiscards)
					metricValue22 := float64(eData.EthernetData.Rt_ifOutErrors)
					metricValue24 := float64(eData.EthernetData.Rt_ifOutMulticastPkts)
					metricValue26 := float64(eData.EthernetData.Rt_ifOutUcastPkts)

						m = append(m, prometheus.MustNewConstMetric(Rt_ifInBroadcastPkts, prometheus.GaugeValue, metricValue3, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifInDiscards, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifInErrors, prometheus.GaugeValue, metricValue5, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifInMulticastPkts, prometheus.GaugeValue, metricValue8, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifInOverSizedPkts, prometheus.GaugeValue, metricValue10, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifInUcastPkts, prometheus.GaugeValue, metricValue11, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifInUndersizedPkts, prometheus.GaugeValue, metricValue12, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifOutBroadcastPkts, prometheus.GaugeValue, metricValue19, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifOutDiscards, prometheus.GaugeValue, metricValue21, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifOutErrors, prometheus.GaugeValue, metricValue22, host.Ip, host.Hostname,"ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifOutMulticastPkts, prometheus.GaugeValue, metricValue24, host.Ip, host.Hostname,"ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifOutUcastPkts, prometheus.GaugeValue, metricValue26, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))

		}
	
	return m
}