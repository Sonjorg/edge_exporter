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
Rt_ifInDiscards		          int    `xml:"rt_ifInDiscards"`
Rt_ifInErrors		          int    `xml:"rt_ifInErrors"`
Rt_ifOutDiscards		      int    `xml:"rt_ifOutDiscards"` //Displays the number of discard errors detected on this port.
Rt_ifOutErrors		          int    `xml:"rt_ifOutErrors"` //Displays the number of errors detected on this port.
}

func EthernetPortCollector(host *config.HostCompose)(m []prometheus.Metric) {

var (
		Rt_ifInDiscards = prometheus.NewDesc("edge_ethernet_ifInDiscards",
			"Displays the number of discard errors detected on this port.",
			[]string{"hostip", "hostname", "job","ethernetportID","ifName","ifAlias"}, nil,
		)
		Rt_ifInErrors = prometheus.NewDesc("edge_ethernet_ifInErrors",
			"Displays the number of errors detected on this port.",
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

					metricValue4 := float64(eData.EthernetData.Rt_ifInDiscards)
					metricValue5 := float64(eData.EthernetData.Rt_ifInErrors)
					metricValue21 := float64(eData.EthernetData.Rt_ifOutDiscards)
					metricValue22 := float64(eData.EthernetData.Rt_ifOutErrors)

						m = append(m, prometheus.MustNewConstMetric(Rt_ifInDiscards, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifInErrors, prometheus.GaugeValue, metricValue5, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifOutDiscards, prometheus.GaugeValue, metricValue21, host.Ip, host.Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						m = append(m, prometheus.MustNewConstMetric(Rt_ifOutErrors, prometheus.GaugeValue, metricValue22, host.Ip, host.Hostname,"ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias))
						

		}
	
	return m
}