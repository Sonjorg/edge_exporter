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

func SystemCollector(host *config.HostCompose) (m []prometheus.Metric, successfulScrape bool) {

	var (
		Rt_CPUUsage = prometheus.NewDesc("edge_system_CPUUsage",
			"Average percent usage of the CPU.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_MemoryUsage = prometheus.NewDesc("edge_system_MemoryUsage",
			"Average percent usage of system memory.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_CPUUptime = prometheus.NewDesc("edge_system_CPUUptime",
			"The total duration in seconds, that the system CPU has been UP and running.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_FDUsage = prometheus.NewDesc("edge_disk_FDUsage",
			"Number of file descriptors used by the system.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_CPULoadAverage1m = prometheus.NewDesc("edge_system_CPULoadAverage1m",
			"Average number of processes over the last one minute waiting to run because CPU is busy.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_CPULoadAverage5m = prometheus.NewDesc("edge_system_CPULoadAverage5m",
			"Average number of processes over the last five minutes waiting to run because CPU is busy.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_CPULoadAverage15m = prometheus.NewDesc("edge_system_CPULoadAverage15m",
			"Average number of processes over the last fifteen minutes waiting to run because CPU is busy.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_TmpPartUsage = prometheus.NewDesc("edge_disk_TmpPartUsage",
			"Percentage of the temporary partition used.",
			[]string{"hostip", "hostname"}, nil,
		)
		Rt_LoggingPartUsage = prometheus.NewDesc("edge_disk_LoggingPartUsage",
			"Percentage of the logging partition used. This is applicable only for the SBC2000.",
			[]string{"hostip", "hostname" }, nil,
		)
		CoreSwitch_Temperature = prometheus.NewDesc("edge_system_temperature",
		"Temperature of the core switch.",
		[]string{"hostip", "hostname" }, nil,
		)
		Error_ip = prometheus.NewDesc("edge_system_status",
			"Returns 1 if the SBC Edge scrape was successful, and 0 if not.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
	)

		dataStr := "https://" + host.Ip + "/rest/system/historicalstatistics/1"
		timeReportedByExternalSystem := time.Now()

		//Trying to fetch chasisinfo from db first, indicated with "null"

		if (!http.SBCIsUp(host.Ip)){ //A quick test to see if contact with sbc
			sbcType, serialNumber, err := utils.GetChassisLabelsDB(host.Ip)
			if err != nil {
				sbcType, serialNumber = "Error fetching chassisinfo", "Error fetching chassisinfo"
				log.Print(err)
			}
			m = append(m, prometheus.NewMetricWithTimestamp(
				timeReportedByExternalSystem,
				prometheus.MustNewConstMetric(
					Error_ip, prometheus.GaugeValue, 0, host.Ip, host.Hostname, sbcType, serialNumber),
			))
			return m, false
		}
		
		phpsessid, err := http.APISessionAuth(host.Username, host.Password, host.Ip)
		if err != nil {
			log.Println("Error retrieving session cookie (system): ", log.Flags(), err)
			sbcType, serialNumber, err := utils.GetChassisLabelsDB(host.Ip)
			if err != nil {
				sbcType, serialNumber = "Error fetching chassisinfo", "Error fetching chassisinfo"
				log.Print(err)
			}
			m = append(m, prometheus.NewMetricWithTimestamp(
				timeReportedByExternalSystem,
				prometheus.MustNewConstMetric(
					Error_ip, prometheus.GaugeValue, 0, host.Ip, host.Hostname, sbcType, serialNumber),
			))

			return m, false//trying next ip address
		}
		//fetching labels from DB or if not exist yet; from router
			sbcType, serialNumber,temperature, err := utils.GetChassisLabelsHTTP(host.Ip, phpsessid)
			if err != nil {
				sbcType, serialNumber,temperature = "Error fetching chassisinfo", "Error fetching chassisinfo",0
				log.Print(err)
			}
		
		//Fetching systemdata
		_, data, err := http.GetAPIData(dataStr, phpsessid)
		if err != nil {
			log.Print("Error collecting from host: ", log.Flags(), err, "\n")
			m = append(m, prometheus.NewMetricWithTimestamp(
				timeReportedByExternalSystem,
				prometheus.MustNewConstMetric(
					Error_ip, prometheus.GaugeValue, 0, host.Ip, host.Hostname, sbcType, serialNumber),
			))
			return m, false
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
		metricValue10 := float64(temperature)

		m = append(m, prometheus.MustNewConstMetric(
			Error_ip, prometheus.GaugeValue, 1, host.Ip, host.Hostname, sbcType, serialNumber))

		m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage15m, prometheus.GaugeValue, metricValue1, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage1m, prometheus.GaugeValue, metricValue2, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage5m, prometheus.GaugeValue, metricValue3, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPUUptime, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPUUsage, prometheus.GaugeValue, metricValue5, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_FDUsage, prometheus.GaugeValue, metricValue6, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_LoggingPartUsage, prometheus.GaugeValue, metricValue7, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_MemoryUsage, prometheus.GaugeValue, metricValue8, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(Rt_TmpPartUsage, prometheus.GaugeValue, metricValue9, host.Ip, host.Hostname))
		m = append(m, prometheus.MustNewConstMetric(CoreSwitch_Temperature, prometheus.GaugeValue, metricValue10, host.Ip, host.Hostname))

	
	return m, true
}
