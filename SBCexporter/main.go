package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
	"log"
	//"bytes"
	//"flag"
	//"io"
	//"log"
	//"net/http"
	//"regexp"
	//"strconv"
	//"time"
	//"github.com/hpcloud/tail"
	//exporter "https://github.com/Sonjorg/HDOmonitoring"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type status struct {
	HTTPcode int `xml:"http_code"`
}

type SBCdata struct {
	XMLname xml.Name `xml:"root"`
	Status  status   `xml:"status"`
	Isdnsg  isdnsg   `xml:"isdnsg"`
	//HTTPcode 				int      `xml:"status>http_code"`
	//Isdnsg					xml.Name `xml:"isdnsg"`
	//ActionsetTableNumber	int		 `xml:"isdnsg>ActionsetTableNumber"`
	//
}

type isdnsg struct {
	Href                  string  `xml:"href,attr"`
	Id                    string  `xml:"id,attr"`
	IncomingCallattempts  int     `xml:"rt_IncomingCallattempts"`
	IncomingCallaccepts   int     `xml:"rt_IncomingCallaccepts"`
}

type fooCollector struct {
	//size     prometheus.Counter
	IncomingCallattempts *prometheus.Desc
	IncomingCallaccepts *prometheus.Desc
	//requests *prometheus.CounterVec
}

func newFooCollector() *fooCollector {
	return &fooCollector{
		IncomingCallattempts: prometheus.NewDesc("incoming_call_attempts",
			"Shows whether a foo has occurred in our cluster",
			nil, nil,
		),
		IncomingCallaccepts: prometheus.NewDesc("incoming_call_accepts",
			"Shows whether a bar has occurred in our cluster",
			nil, nil,
		),
	}
}

	/*func NewConstHistogram(
		desc *Desc,
		count uint64,
		sum float64,
		buckets map[float64]uint64,
		labelValues ...string,
	) (Metric, error)*/

	//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
	func (collector *fooCollector) Describe(ch chan<- *prometheus.Desc) {

		//Update this section with the each metric you create for a given collector
		ch <- collector.IncomingCallaccepts
		ch <- collector.IncomingCallattempts
	}
	//Collect implements required collect function for all promehteus collectors

	func (collector *fooCollector) Collect(ch chan<- prometheus.Metric) {

		//Implement logic here to determine proper metric value to return to prometheus
		//for each descriptor or call other functions that do so.
		var metricValue1 float64

		data, _ := ioutil.ReadFile("APIoutput.xml")
		sbc := &SBCdata{}
		xml.Unmarshal([]byte(data), &sbc)
		fmt.Println("SBC router ID: ",sbc.Isdnsg.Id, "\nRouter href: ",sbc.Isdnsg.Href)
		//var metricValue2 float64
		
		if s, err := strconv.ParseFloat(sbc.Isdnsg.Id, 64); err == nil {
			metricValue1 = s
			fmt.Println(s) // 3.1415927410125732
		}
		

		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		m1 := prometheus.MustNewConstMetric(collector.IncomingCallattempts, prometheus.GaugeValue, metricValue1)
		m2 := prometheus.MustNewConstMetric(collector.IncomingCallaccepts, prometheus.GaugeValue, metricValue1)
		m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
		m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
		ch <- m1
		ch <- m2
	}

func main() {

	foo := newFooCollector()
	prometheus.MustRegister(foo)

	http.Handle("/console/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
	
}
