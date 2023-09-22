/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
package collector

import (
	"encoding/xml"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"log"
	"github.com/prometheus/client_golang/prometheus"
)

// /rest/linecard
type lSBCdata struct {
	XMLname       xml.Name      `xml:"root"`
	LinecardData  LinecardData  `xml:"linecard"`
}
type LinecardData struct {
	Href              string `xml:"href,attr"`
	Rt_CardType		  string `xml:"rt_CardType"`
	Rt_Location		  string `xml:"rt_Location"`
	Rt_ServiceStatus  int    `xml:"rt_ServiceStatus"`
	Rt_Status         int    `xml:"rt_Status"`

}

func LinecardCollector2(host *config.HostCompose)  (m []prometheus.Metric) {

	var (
			Rt_ServiceStatus = prometheus.NewDesc("edge_linecard_ServiceStatus",
				"The service status of the module.",
				[]string{"hostip", "hostname", "job","linecardID","rt_CardType","rt_Location"}, nil,
			)
			Rt_Status = prometheus.NewDesc("edge_linecard_Status",
				"Indicates the hardware initialization state for this card.",
				[]string{"hostip", "hostname", "job","linecardID"}, nil,
			)
		)


		phpsessid,err := http.APISessionAuth(host.Username, host.Password, host.Ip)
		if err != nil {
			log.Print("Error auth", host.Ip, err)
			return 
		}

		//chassis labels from db or http
		chassisType, _, err := utils.GetChassisLabelsDB(host.Ip)
		if err!= nil {
			chassisType = "db chassisData failed"
			log.Print(err)
		}

		var linecardID []string
		// There are two linecard linecardIDs which are different for type of SBC router
		if (chassisType == "SBC1000") {
			linecardID = []string {"7", "8"}
		} else if (chassisType == "SBC2000") {
			linecardID = []string {"1", "2"}
		} else {
			//Couldnt fetch chassis type from db or http: return
			return
		}
			for j := range linecardID {
					url := "https://"+host.Ip+"/rest/linecard/"+linecardID[j]
					_, data, err := http.GetAPIData(url, phpsessid)
						if err != nil {
							log.Print(err)
							continue
						}

					lData := &lSBCdata{}
					err = xml.Unmarshal(data, &lData) //Converting XML data to variables
					if err!= nil {
						log.Print("XML error linecard", err)
						continue
					}
					labelCardType := lData.LinecardData.Rt_CardType
					labelLocation := lData.LinecardData.Rt_Location
					metricValue3 := float64(lData.LinecardData.Rt_ServiceStatus)
					metricValue4 := float64(lData.LinecardData.Rt_Status)
					m = append(m, prometheus.MustNewConstMetric(Rt_ServiceStatus, prometheus.GaugeValue, metricValue3, host.Ip, host.Hostname, "linecard",linecardID[j],labelCardType,labelLocation))
					m = append(m, prometheus.MustNewConstMetric(Rt_Status, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname, "linecard",linecardID[j]))
		}

	return m
}
