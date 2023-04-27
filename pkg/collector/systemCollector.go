package collector

//system status exporter
//rest/system/historicalstatistics/1

import (
	"edge_exporter/pkg/config"
	//"edge_exporter/pkg/database"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"encoding/xml"
	//"fmt"
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

type sMetrics struct {
	Rt_CPUUsage          *prometheus.Desc
	Rt_MemoryUsage       *prometheus.Desc
	Rt_CPUUptime         *prometheus.Desc
	Rt_FDUsage           *prometheus.Desc
	Rt_CPULoadAverage1m  *prometheus.Desc
	Rt_CPULoadAverage5m  *prometheus.Desc
	Rt_CPULoadAverage15m *prometheus.Desc
	Rt_TmpPartUsage      *prometheus.Desc
	Rt_LoggingPartUsage  *prometheus.Desc
	Error_ip             *prometheus.Desc
}

func systemCollector()*sMetrics{

	 return &sMetrics{
		Rt_CPUUsage: prometheus.NewDesc("rt_CPUUsage",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_MemoryUsage: prometheus.NewDesc("rt_MemoryUsage",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_CPUUptime: prometheus.NewDesc("rt_CPUUptime",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_FDUsage: prometheus.NewDesc("rt_FDUsage",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_CPULoadAverage1m: prometheus.NewDesc("rt_CPULoadAverage1m",
			"/rest/system/.",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_CPULoadAverage5m: prometheus.NewDesc("rt_CPULoadAverage5m",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_CPULoadAverage15m: prometheus.NewDesc("rt_CPULoadAverage15m",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_TmpPartUsage: prometheus.NewDesc("Rt_TmpPartUsage",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Rt_LoggingPartUsage: prometheus.NewDesc("Rt_LoggingPartUsage",
			"/rest/system/",
			[]string{"hostip", "hostname", "chassis_type","serial_number"}, nil,
		),
		Error_ip: prometheus.NewDesc("scrape_status",
			"/rest/system/",
			[]string{"hostip", "hostname"}, nil,
		),
	 }
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *sMetrics) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.Rt_CPULoadAverage15m
	ch <- collector.Rt_CPULoadAverage1m
	ch <- collector.Rt_CPULoadAverage5m
	ch <- collector.Rt_CPUUptime
	ch <- collector.Rt_CPUUsage
	ch <- collector.Rt_FDUsage
	ch <- collector.Rt_LoggingPartUsage
	ch <- collector.Rt_MemoryUsage
	ch <- collector.Rt_TmpPartUsage
	ch <- collector.Error_ip

}
//Collect implements required collect function for all promehteus collectors

func (collector *sMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("system")//retrieving targets for this exporter
	if (len(hosts) <= 0) {
		return
	}

	log.Print(hosts)

	for i := 0; i < len(hosts); i++ {
		//authStr := "https://" +hosts[i].ip + "/rest/login"
		dataStr := "https://"+hosts[i].Ip+"/rest/system/historicalstatistics/1"

		//username, password := getAuth(ipaddresses[i])
		timeReportedByExternalSystem := time.Now()//time.Parse(timelayout, mytimevalue)
		chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,"null")
		if err!= nil {
			chassisType, serialNumber = "database failure", "database failure"
			log.Print(err)
		}
		phpsessid,err :=  http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			log.Println("Error retrieving session cookie (system): ",log.Flags(), err)
				 c <- prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, hosts[i].Hostname),
				   )
				   continue //trying next ip address
		}
		//fetching labels from DB or router if not exist yet
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
				  c <- prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, hosts[i].Hostname),
				   )
				continue
		}
		//b := []byte(data) //bytestream
		ssbc := &sSBCdata{}
		err = xml.Unmarshal(data, &ssbc) //Converting XML data to variables
		if err != nil {
			log.Print("XML error system", err)
		}
		//log.Print("Successful API call data: ",ssbc.SystemData)

		metricValue1 := float64(ssbc.SystemData.Rt_CPULoadAverage15m)
		metricValue2 := float64(ssbc.SystemData.Rt_CPULoadAverage1m)
		metricValue3 := float64(ssbc.SystemData.Rt_CPULoadAverage5m)
		metricValue4 := float64(ssbc.SystemData.Rt_CPUUptime)
		metricValue5 := float64(ssbc.SystemData.Rt_CPUUsage)
		metricValue6 := float64(ssbc.SystemData.Rt_FDUsage)
		metricValue7 := float64(ssbc.SystemData.Rt_LoggingPartUsage)
		metricValue8 := float64(ssbc.SystemData.Rt_MemoryUsage)
		metricValue9 := float64(ssbc.SystemData.Rt_TmpPartUsage)

		c <- prometheus.MustNewConstMetric(
			collector.Error_ip, prometheus.GaugeValue, 1, hosts[i].Ip, hosts[i].Hostname)

			c <- prometheus.MustNewConstMetric(collector.Rt_CPULoadAverage15m, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_CPULoadAverage1m, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_CPULoadAverage5m, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_CPUUptime, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_CPUUsage, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_FDUsage, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_LoggingPartUsage, prometheus.GaugeValue, metricValue7, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_MemoryUsage, prometheus.GaugeValue, metricValue8, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_TmpPartUsage, prometheus.GaugeValue, metricValue9, hosts[i].Ip, hosts[i].Hostname,chassisType, serialNumber)
	}
}

/*func sysCollector(collector *sMetrics)  ([]prometheus.Metric) {//(ch chan<- prometheus.Metric){


}*/
// Initializing the exporter
func SystemResourceCollector() {
	hosts := config.GetIncludedHosts("system")//retrieving targets for this exporter
	if (len(hosts) <= 0) {
		return
	}
	c := systemCollector()
	prometheus.MustRegister(c)
}
