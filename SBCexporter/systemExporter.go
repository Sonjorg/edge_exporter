package main

//system status exporter
//rest/system/historicalstatistics/1

import (
	//"crypto/tls"
	//"strings"
	//"bufio"
	"encoding/xml"
	"fmt"
	//"io/ioutil"
	"time"

	//"log"
	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	//"github.com/tiket-oss/phpsessgo"
	//"io/ioutil"
	//"net/http"
	// "net/http/cookiejar"
	// "net/http/cookiejar"
	// "net/url"
	// "regexp"
	// "strconv"
)

type sStatus struct {
	HTTPcode string `xml:"http_code"`
}

type sSBCdata struct {
	XMLname    xml.Name   `xml:"root"`
	Status     sStatus    `xml:"status"`
	SystemData systemData `xml:"historicalstatistics"`

}
type systemData struct {
	Href                 string `xml:"href,attr"`
	Rt_CPUUsage          int    `xml:"rt_CPUUsage"`// Average percent usage of the CPU.
	Rt_MemoryUsage       int    `xml:"rt_MemoryUsage"`// Average percent usage of system memory. int
	Rt_CPUUptime         int    `xml:"rt_CPUUptime"`
	Rt_FDUsage           int    `xml:"rt_FDUsage"`
	Rt_CPULoadAverage1m  int    `xml:"rt_CPULoadAverage1m"`
	Rt_CPULoadAverage5m  int    `xml:"rt_CPULoadAverage5m"`
	Rt_CPULoadAverage15m int    `xml:"rt_CPULoadAverage15m"`
	Rt_TmpPartUsage      int    `xml:"rt_TmpPartUsage"`//Percentage of the temporary partition used. int
	Rt_LoggingPartUsage  int    `xml:"rt_LoggingPartUsage"`
}

type sMetrics struct {
	Href                 *prometheus.Desc
	Rt_CPUUsage          *prometheus.Desc
	Rt_MemoryUsage       *prometheus.Desc
	Rt_CPUUptime         *prometheus.Desc
	Rt_FDUsage           *prometheus.Desc
	Rt_CPULoadAverage1m  *prometheus.Desc
	Rt_CPULoadAverage5m  *prometheus.Desc
	Rt_CPULoadAverage15m *prometheus.Desc
	Rt_TmpPartUsage      *prometheus.Desc
	Rt_LoggingPartUsage  *prometheus.Desc
}

func systemCollector() *sMetrics {
	return &sMetrics{
		Rt_CPUUsage: prometheus.NewDesc("rt_CPUUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_MemoryUsage: prometheus.NewDesc("rt_MemoryUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_CPUUptime: prometheus.NewDesc("rt_CPUUptime",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_FDUsage: prometheus.NewDesc("rt_FDUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_CPULoadAverage1m: prometheus.NewDesc("rt_CPULoadAverage1m",
			"NoDescriptionYet.",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_CPULoadAverage5m: prometheus.NewDesc("rt_CPULoadAverage5m",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_CPULoadAverage15m: prometheus.NewDesc("rt_CPULoadAverage15m",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_TmpPartUsage: prometheus.NewDesc("Rt_TmpPartUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
		),
		Rt_LoggingPartUsage: prometheus.NewDesc("Rt_LoggingPartUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "Href", "HTTP_status"}, nil,
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

}

//Collect implements required collect function for all promehteus collectors

func (collector *sMetrics) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64
	var metricValue4 float64
	var metricValue5 float64
	var metricValue6 float64
	var metricValue7 float64
	var metricValue8 float64
	var metricValue9 float64


	//var HTTPcode float64

	//data, _ := ioutil.ReadFile("sbcsystem.xml")
	phpsessid := APISessionAuth("student", "PanneKake23", "https://10.233.230.11/rest/login")
	data := getAPIData("https://10.233.230.11/rest/system/historicalstatistics/1", phpsessid)
	ssbc := &sSBCdata{}
	b := []byte(data)
	//b, err := ioutil.ReadAll(data)
	/*if err != nil {
	}*/
	xml.Unmarshal([]byte(b), &ssbc)
	fmt.Println(b)

	//fmt.Println(sbc)

	/*if s, err := strconv.ParseFloat(sbc.Status.HTTPcode, 64); err == nil {
		HTTPcode = s
		fmt.Println(s) // 3.1415927410125732
	}*/


	//HTTPcode = float64(sbc.Status.HTTPcode)
	metricValue1 = float64(ssbc.SystemData.Rt_CPULoadAverage15m)
	metricValue2 = float64(ssbc.SystemData.Rt_CPULoadAverage1m)
	metricValue3 = float64(ssbc.SystemData.Rt_CPULoadAverage5m)
	metricValue4 = float64(ssbc.SystemData.Rt_CPUUptime)
	metricValue5 = float64(ssbc.SystemData.Rt_CPUUsage)
	metricValue6 = float64(ssbc.SystemData.Rt_FDUsage)
	metricValue7 = float64(ssbc.SystemData.Rt_LoggingPartUsage)
	metricValue8 = float64(ssbc.SystemData.Rt_MemoryUsage)
	metricValue9 = float64(ssbc.SystemData.Rt_TmpPartUsage)


	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	m1 := prometheus.MustNewConstMetric(collector.Rt_CPULoadAverage15m, prometheus.GaugeValue, metricValue1, "see href","test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)
	m2 := prometheus.MustNewConstMetric(collector.Rt_CPULoadAverage1m, prometheus.GaugeValue, metricValue2, "see href","test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)
	m3 := prometheus.MustNewConstMetric(collector.Rt_CPULoadAverage5m, prometheus.GaugeValue, metricValue3, "see href", "test", "systemstats", ssbc.SystemData.Href,ssbc.Status.HTTPcode)
	m4 := prometheus.MustNewConstMetric(collector.Rt_CPUUptime, prometheus.GaugeValue, metricValue4, "see href", "test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)
	m5 := prometheus.MustNewConstMetric(collector.Rt_CPUUsage, prometheus.GaugeValue, metricValue5, "see href", "test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)
	m6 := prometheus.MustNewConstMetric(collector.Rt_FDUsage, prometheus.GaugeValue, metricValue6, "see href", "test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)
	m7 := prometheus.MustNewConstMetric(collector.Rt_LoggingPartUsage, prometheus.GaugeValue, metricValue7, "see href", "test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)
	m8 := prometheus.MustNewConstMetric(collector.Rt_MemoryUsage, prometheus.GaugeValue, metricValue8, "see href", "test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)
	m9 := prometheus.MustNewConstMetric(collector.Rt_TmpPartUsage, prometheus.GaugeValue, metricValue9, "see href", "test", "systemstats", ssbc.SystemData.Href, ssbc.Status.HTTPcode)

	m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
	m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	m3 = prometheus.NewMetricWithTimestamp(time.Now(), m3)
	m4 = prometheus.NewMetricWithTimestamp(time.Now(), m4)
	m5 = prometheus.NewMetricWithTimestamp(time.Now(), m5)
	m6 = prometheus.NewMetricWithTimestamp(time.Now(), m6)
	m7 = prometheus.NewMetricWithTimestamp(time.Now(), m7)
	m8 = prometheus.NewMetricWithTimestamp(time.Now(), m8)
	m9 = prometheus.NewMetricWithTimestamp(time.Now(), m9)
	ch <- m1
	ch <- m2
	ch <- m3
	ch <- m4
	ch <- m5
	ch <- m6
	ch <- m7
	ch <- m8
	ch <- m9
}


//fetching data from api
func systemExporter() {

	//APISessionAuth()
	/*if err != nil {
		fmt.Println("Apisession auth not working: ", err)
	}*/
	sc := systemCollector()
	prometheus.MustRegister(sc)
	//phpsessid := APISessionAuth()
	//fmt.Println(getAPIData("test", phpsessid))
	//fmt.Println(text)
}


