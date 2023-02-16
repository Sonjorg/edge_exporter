package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	//"bytes"
	//"flag"
	//"io"
	//"log"
	//"net/http"
	//"regexp"
	//"strconv"
	//"time"
	//"github.com/hpcloud/tail"
	//"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
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
	Href   string `xml:"href,attr"`
	Id     string `xml:"id,attr"`
}



/*
	func NewMetrics(reg prometheus.Registerer) *metrics {
		m := &metrics{
			size: prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: "nginx",
				Name:      "size_bytes_total",
				Help:      "Total bytes sent to the clients.",
			}),
			/*requests: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "nginx",
				Name:      "http_requests_total",
				Help:      "Total number of requests.",
			}, []string{"status_code", "method", "path"}),
			duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: "nginx",
				Name:      "http_request_duration_seconds",
				Help:      "Duration of the request.",
				// Optionally configure time buckets
				// Buckets:   prometheus.LinearBuckets(0.01, 0.05, 20),
				Buckets: prometheus.DefBuckets,
			}, []string{"status_code", "method", "path"}),
		}
		reg.MustRegister(m.size)
		return m
	}
*/
func main() {

	/* var (
		targetHost = flag.String("target.host", "localhost", "nginx address with basic_status page")
		targetPort = flag.Int("target.port", 8080, "nginx port with basic_status page")
		targetPath = flag.String("target.path", "/status", "URL path to scrap metrics")
		promPort   = flag.Int("prom.port", 9150, "port to expose prometheus metrics")
		logPath    = flag.String("target.log", "/var/log/nginx/access.log", "path to access.log")
	)
	flag.Parse() */

	data, _ := ioutil.ReadFile("APIoutput.xml")

	sbc := &SBCdata{}

	xml.Unmarshal([]byte(data), &sbc)

	fmt.Println("SBC router ID: ",sbc.Isdnsg.Id, "\nRouter href: ",sbc.Isdnsg.Href)
	//fmt.Println()
	//fmt.Println(note.From)
	//fmt.Println(note.Heading)
	//fmt.Println(note.Body)
}
