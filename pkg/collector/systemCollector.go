/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
package collector

import (
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"encoding/xml"
	"log"
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

type sSBCdata struct {
	XMLname    xml.Name   `xml:"root"`
	Status     sStatus    `xml:"status"`
	SystemData systemData `xml:"historicalstatistics"`
}
type sStatus struct {
	HTTPcode string `xml:"http_code"`
}
type systemData struct {
	Href                 string `xml:"href,attr"`
	Rt_CPUUsage          int    `xml:"rt_CPUUsage"`    // Average percent usage of the CPU.
	Rt_MemoryUsage       int    `xml:"rt_MemoryUsage"` // Average percent usage of system memory. int
	Rt_CPUUptime         int    `xml:"rt_CPUUptime"`
	Rt_FDUsage           int    `xml:"rt_FDUsage"`
	Rt_CPULoadAverage1m  int    `xml:"rt_CPULoadAverage1m"`
	Rt_CPULoadAverage5m  int    `xml:"rt_CPULoadAverage5m"`
	Rt_CPULoadAverage15m int    `xml:"rt_CPULoadAverage15m"`
	Rt_TmpPartUsage      int    `xml:"rt_TmpPartUsage"` //Percentage of the temporary partition used. int
	Rt_LoggingPartUsage  int    `xml:"rt_LoggingPartUsage"`
}

func SystemCollector()(m []prometheus.Metric){

	hosts := config.GetIncludedHosts("system") //retrieving targets for this collector
	if (len(hosts) <= 0) {
		return
	}

	var ( Rt_CPUUsage = prometheus.NewDesc("rt_CPUUsage",
			"Average percent usage of the CPU.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_MemoryUsage = prometheus.NewDesc("rt_MemoryUsage",
			"Average percent usage of system memory.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_CPUUptime = prometheus.NewDesc("rt_CPUUptime",
			"The total duration in seconds, that the system CPU has been UP and running.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_FDUsage = prometheus.NewDesc("rt_FDUsage",
			"Number of file descriptors used by the system.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_CPULoadAverage1m = prometheus.NewDesc("rt_CPULoadAverage1m",
			"Average number of processes over the last one minute waiting to run because CPU is busy.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_CPULoadAverage5m = prometheus.NewDesc("rt_CPULoadAverage5m",
			"Average number of processes over the last five minutes waiting to run because CPU is busy.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_CPULoadAverage15m = prometheus.NewDesc("rt_CPULoadAverage15m",
			"Average number of processes over the last fifteen minutes waiting to run because CPU is busy.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_TmpPartUsage = prometheus.NewDesc("Rt_TmpPartUsage",
			"Percentage of the temporary partition used.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Rt_LoggingPartUsage = prometheus.NewDesc("Rt_LoggingPartUsage",
			"Percentage of the logging partition used. This is applicable only for the SBC2000.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		)
		Error_ip = prometheus.NewDesc("scrape_status",
			"/rest/system/",
			[]string{"hostip", "hostname"}, nil,
		)
	)

	for i := 0; i < len(hosts); i++ {
		dataStr := "https://"+hosts[i].Ip+"/rest/system/historicalstatistics/1"

		timeReportedByExternalSystem := time.Now()
		chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,"null")
		if err!= nil {
			chassisType, serialNumber = "database failure", "database failure"
			log.Print(err)
		}
		phpsessid,err :=  http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			log.Println("Error retrieving session cookie (system): ",log.Flags(), err)
				 m = append(m,  prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, hosts[i].Hostname),
				   ))
				   continue //trying next ip address
		}
		//fetching labels from DB or if not exist yet; from router
		if (chassisType == "database failure") {
			chassisType, serialNumber, err = utils.GetChassisLabels(hosts[i].Ip,phpsessid)
			if err!= nil {
				chassisType, serialNumber = "Failed to get from db", "Failed to get from db"
				log.Print(err)
			}
		}
		//Fetching systemdata
		_, data,err := http.GetAPIData(dataStr, phpsessid)
		if err != nil {
				log.Print("Error collecting from host: ",log.Flags(), err,"\n")
				  m = append(m,  prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, hosts[i].Hostname),
				   ))
				continue
		}
		ssbc := &sSBCdata{}
		err = xml.Unmarshal(data, &ssbc) //Converting XML data to variables
		if err != nil {
			log.Print("XML error system", err)
		}

		metricValue1 := float64(ssbc.SystemData.Rt_CPULoadAverage15m)
		metricValue2 := float64(ssbc.SystemData.Rt_CPULoadAverage1m)
		metricValue3 := float64(ssbc.SystemData.Rt_CPULoadAverage5m)
		metricValue4 := float64(ssbc.SystemData.Rt_CPUUptime)
		metricValue5 := float64(ssbc.SystemData.Rt_CPUUsage)
		metricValue6 := float64(ssbc.SystemData.Rt_FDUsage)
		metricValue7 := float64(ssbc.SystemData.Rt_LoggingPartUsage)
		metricValue8 := float64(ssbc.SystemData.Rt_MemoryUsage)
		metricValue9 := float64(ssbc.SystemData.Rt_TmpPartUsage)

		m = append(m, prometheus.MustNewConstMetric(
			Error_ip, prometheus.GaugeValue, 1, hosts[i].Ip, hosts[i].Hostname))

			m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage15m, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage1m, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage5m, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_CPUUptime, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_CPUUsage, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_FDUsage, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_LoggingPartUsage, prometheus.GaugeValue, metricValue7, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_MemoryUsage, prometheus.GaugeValue, metricValue8, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
			m = append(m, prometheus.MustNewConstMetric(Rt_TmpPartUsage, prometheus.GaugeValue, metricValue9, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber))
	}
	return m
}
