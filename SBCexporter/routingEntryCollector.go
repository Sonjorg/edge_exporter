// routingentry
package main

import (
	"encoding/xml"
	"fmt"
	//"sync"

	//"log"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// rest/routingtable/2/routingentry
// first request
type routingTables struct {
	// Value  float32 `xml:",chardata"`
	XMLName        xml.Name       `xml:"root"`
	RoutingTables2 routingTables2 `xml:"routingtable_list"`
}
type routingTables2 struct {
	RoutingTables3 routingTables3 `xml:"routingtable_pk"`
}
type routingTables3 struct {
	Attr  []string `xml:"id,attr"`
	Value string   `xml:",chardata"`
}

// Second request
type routingEntries struct {
	XMLName       xml.Name      `xml:"root"`
	RoutingEntry2 routingEntry2 `xml:"routingentry_list"`
}
type routingEntry2 struct {
	RoutingEntry3 routingEntry3 `xml:"routingentry_pk"`
}
type routingEntry3 struct {
	Attr  []string `xml:"id,attr"`
	Value string   `xml:",chardata"`
}

// second request
type rSBCdata struct {
	XMLname     xml.Name    `xml:"root"`
	Status      rStatus     `xml:"status"`
	RoutingData routingData `xml:"historicalstatistics"`
}
type rStatus struct {
	HTTPcode string `xml:"http_code"`
}
type routingData struct {
	Href              string `xml:"href,attr"`
	Rt_RuleUsage      int    `xml:"rt_RuleUsage"`
	Rt_ASR            int    `xml:"rt_ASR"`
	Rt_RoundTripDelay int    `xml:"rt_RoundTripDelay"`
	Rt_Jitter         int    `xml:"rt_Jitter"`
	Rt_MOS            int    `xml:"rt_MOS"`
	Rt_QualityFailed  int    `xml:"rt_QualityFailed"`
}

// Metrics for each routingentry
type rMetrics struct {
	Rt_RuleUsage      *prometheus.Desc
	Rt_ASR            *prometheus.Desc
	Rt_RoundTripDelay *prometheus.Desc
	Rt_Jitter         *prometheus.Desc
	Rt_MOS            *prometheus.Desc
	Rt_QualityFailed  *prometheus.Desc
	Error_ip          *prometheus.Desc
}

func routingCollector() *rMetrics {

	return &rMetrics{
		Rt_RuleUsage: prometheus.NewDesc("rt_RuleUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_ASR: prometheus.NewDesc("rt_ASR",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_RoundTripDelay: prometheus.NewDesc("rt_RoundTripDelay",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_Jitter: prometheus.NewDesc("rt_Jitter",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_MOS: prometheus.NewDesc("rt_MOS",
			"NoDescriptionYet.",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_QualityFailed: prometheus.NewDesc("rt_QualityFailed",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_edge_routing",
			"NoDescriptionYet",
			[]string{"Instance", "hostname"}, nil,
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

// Collect implements required collect function for all promehteus collectors
func (collector *rMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := getIncludedHosts("routingentry") //retrieving targets for this exporter
	if len(hosts) <= 0 {
		fmt.Println("no hosts")
		return
	}
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64
	var metricValue4 float64
	var metricValue5 float64
	var metricValue6 float64
	//r := make(map[string][]string)
	//var wg sync.WaitGroup

	for i := range hosts {

			//making map to store ids in sqlite

			phpsessid, err := APISessionAuth(hosts[i].username, hosts[i].password, hosts[i].ip)
			if err != nil {
				fmt.Println("Error auth", hosts[i].ip, err)
				//continue
				return
			}
			_, data, err := getAPIData("https://"+hosts[i].ip+"/rest/routingtable", phpsessid)
			if err != nil {
				fmt.Println("Error routingtable data", hosts[i].ip, err)
				//continue
				return
			}
			//b := []byte(data) //Converting string of data to bytestream
			rt := &routingTables{}
			xml.Unmarshal(data, &rt) //Converting XML data to variables
			//fmt.Println("Successful API call data: ",ssbc.Rt2.Rt3.Attr)
			routingtables := rt.RoutingTables2.RoutingTables3.Attr //ssbc.Rt2.Rt3.Attr
			//fmt.Println("Routingtables " ,routingtables)
			//fmt.Println(b,rt)
			if len(routingtables) <= 0 {
				//return nil, "Routingtables empty"
				fmt.Println("Routingtables empty")
				//continue
				return

			}
			for j := range routingtables {

				//	fmt.Println(phpsessid)
				url := "https://" + hosts[i].ip + "/rest/routingtable/" + routingtables[j] + "/routingentry"
				_, data2, err := getAPIData(url, phpsessid)
				if err != nil {
				}
				//b2 := []byte(data2) //Converting string of data to bytestream
				re := &routingEntries{}
				xml.Unmarshal(data2, &re) //Converting XML data to variables
				routingEntries := re.RoutingEntry2.RoutingEntry3.Attr
				if len(routingEntries) <= 0 {
					fmt.Println("No routingEntry for this routingtable")
					continue
				}
				//
				entries := regexp.MustCompile(`\d+$`)
				var match []string

				for k := range routingEntries {
					tmp := entries.FindStringSubmatch(routingEntries[k])
					match = append(match, tmp)
					//match = entries.FindStringSubmatch(routingEntries[k])
					fmt.Println(tmp)
					fmt.Println(entries.FindStringSubmatch(routingEntries[k]))

				}
				fmt.Println(match)
				//fmt.Println(routingEntries)
				//fmt.Println(match)
				var sqliteDatabase *sql.DB

				sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
				if err != nil {
					fmt.Println(err)
				}
				createRoutingSqlite(sqliteDatabase)
				storeRoutingEntries(sqliteDatabase, hosts[i].ip, "time", routingtables[j], match)
				/*e,err := getRoutingEntries(sqliteDatabase,hosts[i].ip,routingtables[j])
				if err != nil {
					fmt.Println(err)
				}*/
				//fmt.Println("table: ",routingtables[j], "entries: ", e)
				for k := range match {

					url := "https://" + hosts[i].ip + "/rest/routingtable/" + routingtables[j] + "/routingentry/" + match[k] + "/historicalstatistics/1"
					_, data3, err := getAPIData(url, phpsessid)
					if err != nil {
						fmt.Println(err)

						continue
					}

					//fmt.Println(data3)
					//b := []byte(data3) //Converting string of data to bytestream
					rData := &rSBCdata{}
					xml.Unmarshal(data3, &rData) //Converting XML data to variables
					//fmt.Println("Successful API call data: ",rData.RoutingData,"\n")

					metricValue1 = float64(rData.RoutingData.Rt_RuleUsage)
					metricValue2 = float64(rData.RoutingData.Rt_ASR)
					metricValue3 = float64(rData.RoutingData.Rt_RoundTripDelay)
					metricValue4 = float64(rData.RoutingData.Rt_Jitter)
					metricValue5 = float64(rData.RoutingData.Rt_MOS)
					metricValue6 = float64(rData.RoutingData.Rt_QualityFailed)

					c <- prometheus.MustNewConstMetric(collector.Rt_RuleUsage, prometheus.GaugeValue, metricValue1, hosts[i].ip, hosts[i].hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_ASR, prometheus.GaugeValue, metricValue2, hosts[i].ip, hosts[i].hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_RoundTripDelay, prometheus.GaugeValue, metricValue3, hosts[i].ip, hosts[i].hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_Jitter, prometheus.GaugeValue, metricValue4, hosts[i].ip, hosts[i].hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_MOS, prometheus.GaugeValue, metricValue5, hosts[i].ip, hosts[i].hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_QualityFailed, prometheus.GaugeValue, metricValue6, hosts[i].ip, hosts[i].hostname, "routingentry", routingtables[j], match[k], "test")

					//r[routingtables[j]] = append(r[routingtables[j]], match[k])
					//fmt.Println(routingtables[j], match)
					//fmt.Println(r)

				}
				/*var sqliteDatabase *sql.DB

					sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
					if err != nil {
						fmt.Println(err)
					}
					err = createRoutingSqlite(sqliteDatabase)
					if err != nil {
						fmt.Println(err)
					}

					storeRoutingTables(sqliteDatabase, hosts[i].ip, "test", r)

					sqldata, err := getRoutingEntries(sqliteDatabase,hosts[i].ip)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("sqlitedata \n",sqldata)
				}*/
				//if (!routingTablesExists()) {

			}


	}
}

//https://10.233.230.11/rest/routingtable/
//https://10.233.230.11/rest/routingtable/4/routingentry
//https://10.233.230.11/rest/routingtable/2/routingentry/1/historicalstatistics/1

/*func sysCollector(collector *sMetrics)  ([]prometheus.Metric) {//(ch chan<- prometheus.Metric){


}*/
// Initializing the exporter

func routingEntryCollector() {
	c := routingCollector()
	prometheus.MustRegister(c)
}
