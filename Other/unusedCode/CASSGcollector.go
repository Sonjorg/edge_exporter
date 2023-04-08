package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	//"net/http"
	"time"

	//"strconv"
	//"log"
	//"bytes"
	//"flag"
	//"io"

	//"regexp"
	//"strconv"
	//"time"
	//"github.com/hpcloud/tail"
	//exporter "https://github.com/Sonjorg/HDOmonitoring"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

//First request
type cassg struct {
	// Value  float32 `xml:",chardata"`
	 XMLName  xml.Name `xml:"root"`
	 Cassg2   cassg2   `xml:"cassg_list"`
 }
 type cassg2 struct {
	 Cassg3  cassg3 `xml:"_pk"`
 }
 type cassg3 struct {
	 Attr    []string `xml:"id,attr"`
	 Value     string `xml:",chardata"`
 }

//Secnond request
type status struct {
	HTTPcode string `xml:"http_code"`
}

type SBCdata struct {
	XMLname xml.Name `xml:"root"`
	Status  status   `xml:"status"`
	Isdnsg  isdnsg   `xml:"isdnsg"`

}
type isdnsg struct {
	Href                  string `xml:"href,attr"`
	Id                    string `xml:"id,attr"`
    Rt_PeakChannelUsage   int    `xml:"rt_PeakChannelUsage"`
	Rt_CompletedCalls	  int    `xml:"rt_CompletedCalls"`
    Rt_IncompleteCalls    int    `xml:"rt_IncompleteCalls"`
	CustomAdminState      int    `xml:"customAdminState"`
	ApplyToPortName       string `xml:"ApplyToPortName"`
	Description           string `xml:"Description"`
}

//Metrics for each port in use
type imetrics struct {
	Rt_PeakChannelUsage *prometheus.Desc
	Rt_CompletedCalls   *prometheus.Desc
	Rt_IncompleteCalls  *prometheus.Desc
	CustomAdminState    *prometheus.Desc
	//ApplyToPortName     *prometheus.Desc
	//Description         *prometheus.Desc
}

func isdnsgCollector() *imetrics {
	return &imetrics{
		Rt_PeakChannelUsage: prometheus.NewDesc("rt_PeakChannelUsage",
			"Shows incoming call attempts.",
			[]string{"Instance", "hostname", "job", "Id", "ApplyToPortName", "Description", "HTTP_status"}, nil,
		),
		Rt_CompletedCalls: prometheus.NewDesc("rt_CompletedCalls",
			"Shows incoming call attempts.",
			[]string{"Instance", "hostname", "job", "Id", "ApplyToPortName", "Description", "HTTP_status"}, nil,
		),
		Rt_IncompleteCalls: prometheus.NewDesc("rt_IncompleteCalls",
			"Shows incoming call attempts.",
			[]string{"Instance", "hostname", "job", "Id", "ApplyToPortName", "Description", "HTTP_status"}, nil,
		),
		CustomAdminState: prometheus.NewDesc("customAdminState",
			"Shows incoming call attempts.",
			[]string{"Instance", "hostname", "job", "Id", "ApplyToPortName", "Description", "HTTP_status"}, nil,
		),

	}
}


// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *imetrics) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.Rt_PeakChannelUsage
	ch <- collector.Rt_CompletedCalls
	ch <- collector.Rt_IncompleteCalls
}

//Collect implements required collect function for all promehteus collectors

func (collector *metrics) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64

	//var HTTPcode float64

	data, _ := ioutil.ReadFile("APIoutput.xml")
	sbc := &SBCdata{}
	xml.Unmarshal([]byte(data), &sbc)
	fmt.Println("Incoming call attempts / accepts: ", sbc.Isdnsg.IncomingCallattempts, "/", sbc.Isdnsg.IncomingCallaccepts, "\nSBC ID: ", sbc.Isdnsg.Id, "\nRouter href: ", sbc.Isdnsg.Href)

	/*if s, err := strconv.ParseFloat(sbc.Status.HTTPcode, 64); err == nil {
		HTTPcode = s
		fmt.Println(s) // 3.1415927410125732
	}*/
	//HTTPcode = float64(sbc.Status.HTTPcode)
	metricValue1 = float64(sbc.Isdnsg.IncomingCallattempts)
	metricValue2 = float64(sbc.Isdnsg.IncomingCallaccepts)

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	m1 := prometheus.MustNewConstMetric(collector.IncomingCallattempts, prometheus.GaugeValue, metricValue1, sbc.Isdnsg.Id, sbc.Isdnsg.Href,sbc.Status.HTTPcode)
	m2 := prometheus.MustNewConstMetric(collector.IncomingCallaccepts, prometheus.GaugeValue, metricValue2, sbc.Isdnsg.Id, sbc.Isdnsg.Href,sbc.Status.HTTPcode)
	m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
	m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	ch <- m1
	ch <- m2
}