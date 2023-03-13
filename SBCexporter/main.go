package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	//"strconv"
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
	IncomingCallattempts  int    `xml:"rt_IncomingCallattempts"`
	IncomingCallaccepts   int    `xml:"rt_IncomingCallaccepts"`
	IncomingCallrejects   int    `xml:"rt_IncomingCallrejects"`
	IncomingCallcompletes int    `xml:"rt_IncomingCallcompletes"`
	OutgoingCallattempts  int    `xml:"rt_OutgoingCallattempts"`
	OutgoingCallaccepts   int    `xml:"rt_OutgoingCallaccepts"`
	OutgoingCallrejects   int    `xml:"rt_OutgoingCallrejects"`
	OutgoingCallcompletes int    `xml:"rt_OutgoingCallcompletes"`
}

type metrics struct {
	//size     prometheus.Counter
	IncomingCallattempts *prometheus.Desc
	IncomingCallaccepts  *prometheus.Desc
	//requests *prometheus.CounterVec
}

func newFooCollector() *metrics {
	return &metrics{
		IncomingCallattempts: prometheus.NewDesc("incoming_call_attempts",
			"Shows incoming call attempts.",
			[]string{"Id", "Href", "HTTP_status"}, nil,
		),
		IncomingCallaccepts: prometheus.NewDesc("incoming_call_accepts",
			"Shows incoming call accepts.",
			[]string{"Id", "Href","HTTP_status"}, nil,
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

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *metrics) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.IncomingCallaccepts
	ch <- collector.IncomingCallattempts
}

//Collect implements required collect function for all promehteus collectors

func (collector *metrics) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue1 float64
	var metricValue2 float64
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

func main() {

	if (len(getIpAdrExp("systemstats")) > 0) {
		systemDataExporter()
	}
	/*if (len(getIpAdrExp("callstats")) > 0) {
		callDataExporter()
	}*/
	foo := newFooCollector()
	prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9100", nil))

}
