// routingentry
package collector

import (
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"encoding/xml"
	"fmt"
	//"sync"
	//"log"
	"regexp"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
	//"exporter/sqlite"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

)

// first request
// rest/routingtable/
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
// rest/routingtable/4/routingentry
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

// Third request
// rest/routingtable/2/routingentry/1/historicalstatistics/1
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
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "chassis_type","serial_number"}, nil,
		),
		Rt_ASR: prometheus.NewDesc("rt_ASR",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "chassis_type","serial_number"}, nil,
		),
		Rt_RoundTripDelay: prometheus.NewDesc("rt_RoundTripDelay",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "chassis_type","serial_number"}, nil,
		),
		Rt_Jitter: prometheus.NewDesc("rt_Jitter",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "chassis_type","serial_number"}, nil,
		),
		Rt_MOS: prometheus.NewDesc("rt_MOS",
			"NoDescriptionYet.",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "chassis_type","serial_number"}, nil,
		),
		Rt_QualityFailed: prometheus.NewDesc("rt_QualityFailed",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "chassis_type","serial_number"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_edge_routing",
			"NoDescriptionYet",
			[]string{"Instance", "hostname","job","routing_table", "error_reason","chassis_type","serial_number"}, nil,
		),
	}
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
	hosts := config.GetIncludedHosts("routingentry") //retrieving targets for this exporter
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

	//var timeLast string
	var sqliteDatabase *sql.DB
	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		fmt.Println(err)
	}
	for i := range hosts {
		var timeLastString string //Fetched from database, the routingentries and tables are stored for 24 hours as requested by HDO
		var timeLast time.Time
		phpsessid, err := http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
			if err != nil {
				fmt.Println("Error authentication", hosts[i].Ip, err)
				c <- prometheus.MustNewConstMetric(
					collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, "routingentry","NA", "authentication error")
				continue
			}
			var routingtables []string
			var routingEntryMap = make(map[string][]string)
			var DBexists bool = database.RoutingTablesExists(sqliteDatabase,hosts[i].Ip) //Previous data is stored in db? Fetch this data
			//fmt.Println("exists:",exists)
			if (DBexists) {
				routingEntryMap,routingtables,timeLastString,err = database.GetRoutingData(sqliteDatabase,hosts[i].Ip) // From db: returning a map of routingentables to routingentries (array),
				if err != nil {
					fmt.Println(err)
				}
				timeLast,err = time.Parse(time.RFC3339, timeLastString)
				if err != nil {
					fmt.Println(err)
				}
			}
			//If 24 hours has not passed since last data was stored in database, use this data
			//fmt.Println(b)
			if (!DBexists || database.TimeIsUp(24, timeLast) == true)  { //Routing data has expired, fetching new routingentries
				_, data, err := http.GetAPIData("https://"+hosts[i].Ip+"/rest/routingtable", phpsessid)
				if err != nil {
					fmt.Println("Error routingtable data", hosts[i].Ip, err)
					c <- prometheus.MustNewConstMetric(
							collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, "routingentry","NA", "no_routing_tables")

					continue
				}
				rt := &routingTables{}
				err = xml.Unmarshal(data, &rt) //Converting XML data to variables
				if err != nil {
					fmt.Println("XML Conversion error", err)
				}
				fmt.Println("fetched from router")
				routingtables = rt.RoutingTables2.RoutingTables3.Attr
			}
							//using previous routingentries if within time

			if len(routingtables) <= 0 {
				fmt.Println("Routingtables empty")
				continue //routingtables emtpy, try next host
			}

			chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,phpsessid)
			if err!= nil {
				chassisType, serialNumber = "database failure", "database failure"
				fmt.Println(err)
			}

			for j := range routingtables {
				var match []string //variable to hold routingentries cleaned with regex
				//Trying to fetch routingentries from database, if not exist yet, fetch new ones
				if (DBexists) {
					for k,v := range routingEntryMap {// fetching routingentries from a map from db
						if (k == routingtables[j]) {
							for re := range v {
								match = append(match,v[re]) //using previous routingentries (match)
							}
						}
					}

				} else { // DB doesn't exist, so fetch new routingentries with
					url := "https://" + hosts[i].Ip + "/rest/routingtable/" + routingtables[j] + "/routingentry"
					_, data2, err := http.GetAPIData(url, phpsessid)
					if err != nil {
					}
					//b2 := []byte(data2) //Converting string of data to bytestream
					re := &routingEntries{}
					xml.Unmarshal(data2, &re) //Converting XML data to variables
					routingEntries := re.RoutingEntry2.RoutingEntry3.Attr
					/*if len(routingEntries) <= 0 {
						c <- prometheus.MustNewConstMetric(
							collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip, "systemstats", "no_routing_entries")
							continue

					}*/
					//
					entries := regexp.MustCompile(`\d+$`)

					//Because routingentries from the hosts are displayed as a list of for example "2:4", "2:5", we are using regex to get only the routingentries
					for k := range routingEntries {
						tmp := entries.FindStringSubmatch(routingEntries[k])
						for l := range tmp {
							match = append(match, tmp[l])
							//fmt.Println(tmp[l])
						}
					}
					//Storing fetched routingentries
					/*err = database.CreateRoutingSqlite(sqliteDatabase)
					if err != nil {
						fmt.Println(err)
					}*/
					now := time.Now().Format(time.RFC3339)
					fmt.Println("NOW:", now)
					err = database.StoreRoutingEntries(sqliteDatabase, hosts[i].Ip, now, routingtables[j], match)
					if err != nil {
						fmt.Println(err)
					}
				}

				if (len(match) <= 0) {
						continue
				}
				for k := range match {

					url := "https://" + hosts[i].Ip + "/rest/routingtable/" + routingtables[j] + "/routingentry/" + match[k] + "/historicalstatistics/1"
					_, data3, err := http.GetAPIData(url, phpsessid)
					if err != nil {
						fmt.Println(err)

						continue
					}

					rData := &rSBCdata{}
					xml.Unmarshal(data3, &rData) //Converting XML data to variables
					if err!= nil {
						fmt.Println("XML error routing", err)
						//continue
					}
					//fmt.Println("Successful API call data: ",rData.RoutingData)

					metricValue1 = float64(rData.RoutingData.Rt_RuleUsage)
					metricValue2 = float64(rData.RoutingData.Rt_ASR)
					metricValue3 = float64(rData.RoutingData.Rt_RoundTripDelay)
					metricValue4 = float64(rData.RoutingData.Rt_Jitter)
					metricValue5 = float64(rData.RoutingData.Rt_MOS)
					metricValue6 = float64(rData.RoutingData.Rt_QualityFailed)

					c <- prometheus.MustNewConstMetric(collector.Rt_RuleUsage, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], chassisType,serialNumber)
					c <- prometheus.MustNewConstMetric(collector.Rt_ASR, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], chassisType,serialNumber)
					c <- prometheus.MustNewConstMetric(collector.Rt_RoundTripDelay, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], chassisType,serialNumber)
					c <- prometheus.MustNewConstMetric(collector.Rt_Jitter, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], chassisType,serialNumber)
					c <- prometheus.MustNewConstMetric(collector.Rt_MOS, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], chassisType,serialNumber)
					c <- prometheus.MustNewConstMetric(collector.Rt_QualityFailed, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], chassisType,serialNumber)
				}
			}
	}
}
/*
func RoutingTablesExists(sqliteDatabase *sql.DB, s string) {
	panic("unimplemented")
}*/

func RoutingEntryCollector() {
	c := routingCollector()
	prometheus.MustRegister(c)
}
