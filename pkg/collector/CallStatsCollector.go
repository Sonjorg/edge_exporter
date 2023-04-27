package collector

//system status exporter
//rest/system/historicalstatistics/1

import (
	"encoding/xml"
	//"fmt"
	"log"
	"github.com/prometheus/client_golang/prometheus"
	//"time"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/config"
	//"edge_exporter/pkg/utils"

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
	//Error_ip                    *prometheus.Desc
}

func callStats()*cMetrics{

	 return &cMetrics{
		Rt_NumCallAttempts: prometheus.NewDesc("rt_NumCallAttempts",
			"systemcallstats",
			[]string{"hostip", "hostname"}, nil,
		),
		Rt_NumCallSucceeded: prometheus.NewDesc("rt_NumCallSucceeded",
			"systemcallstats",
			[]string{"hostip", "hostname"}, nil,
		),
		Rt_NumCallFailed: prometheus.NewDesc("rt_NumCallFailed",
			"systemcallstats",
			[]string{"hostip", "hostname"}, nil,
		),
		Rt_NumCallCurrentlyUp: prometheus.NewDesc("rt_NumCallCurrentlyUp",
			"systemcallstats",
			[]string{"hostip", "hostname"}, nil,
		),
		Rt_NumCallAbandonedNoTrunk: prometheus.NewDesc("rt_NumCallAbandonedNoTrunk",
			"systemcallstats.",
			[]string{"hostip", "hostname"}, nil,
		),
		Rt_NumCallUnAnswered: prometheus.NewDesc("rt_NumCallUnAnswered",
			"systemcallstats",
			[]string{"hostip", "hostname"}, nil,
		),
		/*Error_ip: prometheus.NewDesc("error_edge_call_stats",
			"systemcallstats",
			[]string{"hostip", "hostname", "Href","chassis_type","serial_number"}, nil,
		),*/
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
	//ch <- collector.Error_ip

}
//Collect implements required collect function for all promehteus collectors

func (collector *cMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("systemcallstats") //retrieving targets for this collector
	if (len(hosts) <= 0) {
		return
	}

	for i := 0; i < len(hosts); i++ {

		phpsessid,err :=  http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			log.Print("Error retrieving session cookie: ", err,"\n")
				   continue //trying next ip address
		}

		dataStr := "https://"+hosts[i].Ip+"/rest/systemcallstats"
		_, data,err := http.GetAPIData(dataStr, phpsessid)
		if err != nil {
				log.Print("Error collecting systemcall data: ", err,"\n")
				continue
		}

		b := []byte(data) //Converting string of data to bytestream
		ssbc := &cSBCdata{}
		err = xml.Unmarshal(b, &ssbc) //Converting XML data to variables
		if err!= nil {
			log.Print("XML error callstats", err)
			continue
		}

		metricValue1 := float64(ssbc.CallStatsData.Rt_NumCallAttempts)
		metricValue2 := float64(ssbc.CallStatsData.Rt_NumCallSucceeded)
		metricValue3 := float64(ssbc.CallStatsData.Rt_NumCallFailed)
		metricValue4 := float64(ssbc.CallStatsData.Rt_NumCallCurrentlyUp)
		metricValue5 := float64(ssbc.CallStatsData.Rt_NumCallAbandonedNoTrunk)
		metricValue6 := float64(ssbc.CallStatsData.Rt_NumCallUnAnswered)
		// ssbc.CallStatsData.Href
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallAttempts, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallSucceeded, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallFailed, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallCurrentlyUp, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallAbandonedNoTrunk, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname)
			c <- prometheus.MustNewConstMetric(collector.Rt_NumCallUnAnswered, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname)
	}
}

// Initializing the exporter
func CallStatsCollector() {
	hosts := config.GetIncludedHosts("systemcallstats") //retrieving targets for this collector
	if (len(hosts) <= 0) {
		return
	}
	c := callStats()
	prometheus.MustRegister(c)
}
