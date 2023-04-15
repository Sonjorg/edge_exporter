package collector

//system status exporter
//rest/system/historicalstatistics/1

import (
	"encoding/xml"
	"fmt"
	"log"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/utils"

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
	Href                     string `xml:"href,attr"`
	Rt_NumCallAttempts          int `xml:"rt_NumCallAttempts"`        // Average percent usage of the CPU.
	Rt_NumCallSucceeded         int `xml:"rt_NumCallSucceeded"`       // Average percent usage of system memory. int
	Rt_NumCallFailed            int `xml:"rt_NumCallFailed"`
	Rt_NumCallCurrentlyUp       int `xml:"rt_NumCallCurrentlyUp"`     //Number of currently connected calls system wide. int
	Rt_NumCallAbandonedNoTrunk  int `xml:"rt_NumCallAbandonedNoTrunk"`//Number of rejected calls due to no channel available system wide since system came up.
	Rt_NumCallUnAnswered        int `xml:"rt_NumCallUnAnswered"`
}

type cMetrics struct {
	Rt_NumCallAttempts          *prometheus.Desc
	Rt_NumCallSucceeded         *prometheus.Desc
	Rt_NumCallFailed            *prometheus.Desc
	Rt_NumCallCurrentlyUp       *prometheus.Desc
	Rt_NumCallAbandonedNoTrunk  *prometheus.Desc
	Rt_NumCallUnAnswered        *prometheus.Desc
	Error_ip                    *prometheus.Desc
}

func callStats()*cMetrics{

	 return &cMetrics{
		Rt_NumCallAttempts: prometheus.NewDesc("rt_NumCallAttempts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "chassis_type","serial_number"}, nil,
		),
		Rt_NumCallSucceeded: prometheus.NewDesc("rt_NumCallSucceeded",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "chassis_type","serial_number"}, nil,
		),
		Rt_NumCallFailed: prometheus.NewDesc("rt_NumCallFailed",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "chassis_type","serial_number"}, nil,
		),
		Rt_NumCallCurrentlyUp: prometheus.NewDesc("rt_NumCallCurrentlyUp",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "chassis_type","serial_number"}, nil,
		),
		Rt_NumCallAbandonedNoTrunk: prometheus.NewDesc("rt_NumCallAbandonedNoTrunk",
			"NoDescriptionYet.",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "chassis_type","serial_number"}, nil,
		),
		Rt_NumCallUnAnswered: prometheus.NewDesc("rt_NumCallUnAnswered",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "chassis_type","serial_number"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_edge_call_stats",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href","chassis_type","serial_number"}, nil,
		),
	 }
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *cMetrics) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.Rt_NumCallAttempts
	ch <- collector.Rt_NumCallSucceeded
	ch <- collector.Rt_NumCallFailed
	ch <- collector.Rt_NumCallCurrentlyUp
	ch <- collector.Rt_NumCallAbandonedNoTrunk
	ch <- collector.Rt_NumCallUnAnswered
	ch <- collector.Error_ip

}
//Collect implements required collect function for all promehteus collectors

func (collector *cMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("systemcallstats") //retrieving targets for this collector
	if (len(hosts) <= 0) {
		return
	}
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64
	var metricValue4 float64
	var metricValue5 float64
	var metricValue6 float64

	//fmt.Println(hosts)

	for i := 0; i < len(hosts); i++ {
		nr := strconv.Itoa(i)
		//authStr := "https://" +hosts[i].ip + "/rest/login"

		//username, password := getAuth(ipaddresses[i])
		timeReportedByExternalSystem := time.Now()//time.Parse(timelayout, mytimevalue)
		phpsessid,err :=  http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			log.Println("Error retrieving session cookie: ",log.Flags(), err,"\n")
			//return nil, err <-this line would result in error for systemexp on all hosts
			//returning a prometheus error metric
				 c <- prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, hosts[i].Hostname, "systemcallstats",nr, "/rest/systemcallstats"),
				   )
				   continue //trying next ip address
		}
		//chassis labels from db or http
		chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,phpsessid)
		if err!= nil {
			chassisType, serialNumber = "db chassisData fail", "db chassisData fail"
			fmt.Println(err)
		}
		dataStr := "https://"+hosts[i].Ip+"/rest/systemcallstats"
		_, data,err := http.GetAPIData(dataStr, phpsessid)
		if err != nil {
				fmt.Println("Error collecting from host: ",log.Flags(), err,"\n")
				  c <- prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, hosts[i].Hostname, "systemcallstats",nr, "/rest/systemcallstats",chassisType, serialNumber),
				   )
				continue
		}
		b := []byte(data) //Converting string of data to bytestream
		ssbc := &cSBCdata{}
		err = xml.Unmarshal(b, &ssbc) //Converting XML data to variables
		if err!= nil {
			fmt.Println("XML error callstats", err)
			//continue
		}
		fmt.Println("Successful API call data: ", ssbc.CallStatsData)

		metricValue1 = float64(ssbc.CallStatsData.Rt_NumCallAttempts)
		metricValue2 = float64(ssbc.CallStatsData.Rt_NumCallSucceeded)
		metricValue3 = float64(ssbc.CallStatsData.Rt_NumCallFailed)
		metricValue4 = float64(ssbc.CallStatsData.Rt_NumCallCurrentlyUp)
		metricValue5 = float64(ssbc.CallStatsData.Rt_NumCallAbandonedNoTrunk)
		metricValue6 = float64(ssbc.CallStatsData.Rt_NumCallUnAnswered)

			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallAttempts, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname, "systemstats",nr, ssbc.CallStatsData.Href, chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallSucceeded, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname, "systemstats",nr, ssbc.CallStatsData.Href, chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallFailed, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname, "systemstats",nr, ssbc.CallStatsData.Href, chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallCurrentlyUp, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname, "systemstats",nr, ssbc.CallStatsData.Href, chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallAbandonedNoTrunk, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname, "systemstats",nr, ssbc.CallStatsData.Href, chassisType, serialNumber)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallUnAnswered, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname, "systemstats",nr, ssbc.CallStatsData.Href, chassisType, serialNumber)
	}
}

// Initializing the exporter
func CallStatsCollector() {
		c := callStats()
		prometheus.MustRegister(c)
}
