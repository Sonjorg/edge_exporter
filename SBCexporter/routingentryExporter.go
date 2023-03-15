//routingentry
package main

//system status exporter
//rest/system/historicalstatistics/1

import (
	"encoding/xml"
	"fmt"
	"log"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

//rest/routingtable/2/routingentry
//first request
type call1xml1 struct {
	// Value  float32 `xml:",chardata"`
	 XMLName    xml.Name  `xml:"root"`
	 Call1xml2  call1xml2 `xml:"routingtable_list"`
 }
 type call1xml2 struct {
	 Call1xml2  call1xml3 `xml:"routingtable_pk"`
 }
 type call1xml3 struct {
	 Attr    []string `xml:"id,attr"`
	 Value   string `xml:",chardata"`

 }
//https://10.233.230.11/rest/routingtable/2/routingentry/ + ssbc.Rt2.Rt3.Attr[j]

//second request
type rSBCdata struct {
	XMLname    xml.Name   `xml:"root"`
	Status     rStatus    `xml:"status"`
	routingData routingData `xml:"historicalstatistics"`
}
type rStatus struct {
	HTTPcode string `xml:"http_code"`
}
type routingData struct {
Href                string `xml:"href,attr"`
Rt_RuleUsage		int    `xml:"rt_RuleUsage"`
Rt_ASR				int    `xml:"rt_ASR"`
Rt_RoundTripDelay	int    `xml:"rt_RoundTripDelay"`
Rt_Jitter           int    `xml:"rt_Jitter"`
Rt_MOS              int    `xml:"rt_MOS"`
Rt_QualityFailed    int    `xml:"rt_QualityFailed"`
}

type rMetrics struct {
	Rt_RuleUsage		*prometheus.Desc
	Rt_ASR				*prometheus.Desc
	Rt_RoundTripDelay	*prometheus.Desc
	Rt_Jitter           *prometheus.Desc
	Rt_MOS              *prometheus.Desc
	Rt_QualityFailed    *prometheus.Desc
	Error_ip            *prometheus.Desc
}

func routingCollector()*rMetrics{

	 return &rMetrics{
		Rt_RuleUsage: prometheus.NewDesc("rt_RuleUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "HTTP_status"}, nil,
		),
		Rt_ASR: prometheus.NewDesc("rt_ASR",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "HTTP_status"}, nil,
		),
		Rt_RoundTripDelay: prometheus.NewDesc("rt_RoundTripDelay",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "HTTP_status"}, nil,
		),
		Rt_Jitter: prometheus.NewDesc("rt_Jitter",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "HTTP_status"}, nil,
		),
		Rt_MOS: prometheus.NewDesc("rt_MOS",
			"NoDescriptionYet.",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "HTTP_status"}, nil,
		),
		Rt_QualityFailed: prometheus.NewDesc("rt_QualityFailed",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "HTTP_status"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_edge_routing",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","host_nr", "Href", "HTTP_status"}, nil,
		),
	 }

//}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *rMetrics) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.Rt_RuleUsage
	ch <- collector.Rt_ASR
	ch <- collector.Rt_RoundTripDelay
	ch <- collector.Rt_Jitter
	ch <- collector.Rt_MOS
	ch <- collector.Rt_QualityFailed
	ch <- collector.Error_ip
}
//Collect implements required collect function for all promehteus collectors
func (collector *rMetrics) Collect(c chan<- prometheus.Metric) {
	//metrics := sysCollector(collector) //NB: Errors are returned as array of NewInvalidMetric()
	//array of metrics is sent through the channel
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64
	var metricValue4 float64
	var metricValue5 float64
	var metricValue6 float64
	//var metricsValuex float64

	var username string
	var password string
	//var phpsessid string
	ipaddresses := getIpAdrExp("systemStats") //retrieving sources for this exporter

	//DO NOT DELETE: ipaddresses = getIPNotExl("systemExporter", testConfig)
	//phpsessid := APISessionAuth("student", "PanneKake23", "https://10.233.230.11/rest/login")
	username = "student"
	password = "PanneKake23"
	//var err error
	fmt.Println(ipaddresses)


	for i := 0; i < len(ipaddresses); i++ {
		nr := strconv.Itoa(i)
		username = `student`
		password = `PanneKake23`
		authStr := "https://" +ipaddresses[i] + "/rest/login"
		//m := make(map[string]string)
		//m["route"] = "egegrreg"
		fmt.Println("Api call on ip: ",ipaddresses[i],"\n")

		timeReportedByExternalSystem := time.Now()//time.Parse(timelayout, mytimevalue)
		phpsessid,err :=  APISessionAuth(username, password,authStr)
		if err != nil {
			log.Println("Error retrieving session cookie: ",log.Flags(), err,"\n")
			//return nil, err <-this line would result in error for systemexp on all hosts
			//returning a prometheus error metric
				 c <- prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						collector.Error_ip, prometheus.GaugeValue, 0, ipaddresses[i], "test", "systemstats-host-"+ipaddresses[i],nr, ipaddresses[i], "500"),
				   )
				   continue //trying next ip address
		}
		//https://10.233.230.11/rest/routingtable/
		//https://10.233.230.11/rest/routingtable/4/routingentry
		//https://10.233.230.11/rest/routingtable/2/routingentry/1/historicalstatistics/1

		dataStr := "https://"+ipaddresses[i]+"/rest/routingentry/historicalstatistics/1"
		data,err := getAPIData(dataStr, phpsessid)
		if err != nil {
				fmt.Println("Error collecting from host: ",log.Flags(), err,"\n")
				  c <- prometheus.NewMetricWithTimestamp(
					timeReportedByExternalSystem,
					prometheus.MustNewConstMetric(
						collector.Error_ip, prometheus.GaugeValue, 0, ipaddresses[i], "test", "systemstats-host-"+ipaddresses[i],nr, ipaddresses[i], "55"),
				   )
				continue
		}
		b := []byte(data) //Converting string of data to bytestream
		ssbc := &rSBCdata{}
		xml.Unmarshal(b, &ssbc) //Converting XML data to variables
		fmt.Println("Successful API call data: ",ssbc.routingData,"\n")

		metricValue1 = float64(ssbc.routingData.Rt_RuleUsage)
		metricValue2 = float64(ssbc.routingData.Rt_ASR)
		metricValue3 = float64(ssbc.routingData.Rt_RoundTripDelay)
		metricValue4 = float64(ssbc.routingData.Rt_Jitter)
		metricValue5 = float64(ssbc.routingData.Rt_MOS)
		metricValue6 = float64(ssbc.routingData.Rt_QualityFailed)

			c <- prometheus.MustNewConstMetric(collector.Rt_RuleUsage, prometheus.GaugeValue, metricValue1, ipaddresses[i], "test", "routingentry",nr, ssbc.routingData.Href, ssbc.Status.HTTPcode)
			c <- prometheus.MustNewConstMetric(collector.Rt_ASR, prometheus.GaugeValue, metricValue2, ipaddresses[i], "test", "routingentry",nr, ssbc.routingData.Href, ssbc.Status.HTTPcode)
			c <- prometheus.MustNewConstMetric(collector.Rt_RoundTripDelay, prometheus.GaugeValue, metricValue3, ipaddresses[i], "test", "routingentry",nr, ssbc.routingData.Href, ssbc.Status.HTTPcode)
			c <- prometheus.MustNewConstMetric(collector.Rt_Jitter, prometheus.GaugeValue, metricValue4, ipaddresses[i], "test", "routingentry",nr, ssbc.routingData.Href, ssbc.Status.HTTPcode)
			c <- prometheus.MustNewConstMetric(collector.Rt_MOS, prometheus.GaugeValue, metricValue5, ipaddresses[i], "test", "routingentry",nr, ssbc.routingData.Href, ssbc.Status.HTTPcode)
			c <- prometheus.MustNewConstMetric(collector.Rt_QualityFailed, prometheus.GaugeValue, metricValue6, ipaddresses[i], "test", "routingentry",nr, ssbc.routingData.Href, ssbc.Status.HTTPcode)
	}
}

/*func sysCollector(collector *sMetrics)  ([]prometheus.Metric) {//(ch chan<- prometheus.Metric){


}*/
// Initializing the exporter
func test() {
		sc := routingCollector()
		prometheus.MustRegister(sc)
}
