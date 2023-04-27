// routingentry
package collector

import (
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"encoding/xml"
	"fmt"

	//"fmt"
	//"sync"
	"log"
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
			"routingentry",
			[]string{"hostip", "hostname",  "routing_table", "routing_entry"}, nil,
		),
		Rt_ASR: prometheus.NewDesc("rt_ASR",
			"routingentry",
			[]string{"hostip", "hostname",  "routing_table", "routing_entry"}, nil,
		),
		Rt_RoundTripDelay: prometheus.NewDesc("rt_RoundTripDelay",
			"routingentry",
			[]string{"hostip", "hostname",  "routing_table", "routing_entry"}, nil,
		),
		Rt_Jitter: prometheus.NewDesc("rt_Jitter",
			"routingentry",
			[]string{"hostip", "hostname",  "routing_table", "routing_entry"}, nil,
		),
		Rt_MOS: prometheus.NewDesc("rt_MOS",
			"routingentry.",
			[]string{"hostip", "hostname",  "routing_table", "routing_entry"}, nil,
		),
		Rt_QualityFailed: prometheus.NewDesc("rt_QualityFailed",
			"routingentry",
			[]string{"hostip", "hostname",  "routing_table", "routing_entry"}, nil,
		),
		/*Error_ip: prometheus.NewDesc("error_edge_routing",
			"routingentry",
			[]string{"hostip", "hostname","routing_table", "error_reason","chassis_type","serial_number"}, nil,
		),*/
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
	//ch <- collector.Error_ip
}

// Collect implements required collect function for all promehteus collectors
func (collector *rMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("routingentry") //retrieving targets for this exporter
	if len(hosts) <= 0 {
		log.Print("no hosts")
		return
	}


	//var timeLast string
	var sqliteDatabase *sql.DB
	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Print(err)
	}
	for i := range hosts {
		var timeLastString string //Fetched from database, the routingentries and tables are stored for 24 hours as requested by HDO
		var timeLast time.Time
		phpsessid, err := http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
			if err != nil {
				log.Print("Error authentication", hosts[i].Ip, err)
				continue
			}
			var routingtables []string
			var routingEntryMap = make(map[string][]string)
			var DBexists bool = database.RoutingTablesExists(sqliteDatabase,hosts[i].Ip) //Previous data is stored in db? Fetch this data
			if (DBexists) {
				routingEntryMap,routingtables,timeLastString,err = database.GetRoutingData(sqliteDatabase,hosts[i].Ip) // From db: returning a map of routingentables to routingentries (array),
				if err != nil {
					log.Print(err)
				}
				timeLast,err = time.Parse(time.RFC3339, timeLastString)
				if err != nil {
					log.Print(err)
				}
			}
			timeSchedule := hosts[i].RoutingEntryTime
			//If 24 hours has not passed since last data was stored in database, use this data
			if (!DBexists || utils.Expired(timeSchedule, timeLast))  { //Routing data has expired, fetching new routingentries
				fmt.Println("Fetching routing data from http")
				_, data, err := http.GetAPIData("https://"+hosts[i].Ip+"/rest/routingtable", phpsessid)
				if err != nil {
					log.Print("Error routingtable data", hosts[i].Ip, err)
					continue
				}
				rt := &routingTables{}
				err = xml.Unmarshal(data, &rt) //Converting XML data to variables
				if err != nil {
					log.Print("XML error routingentry", err)
					continue
				}
				routingtables = rt.RoutingTables2.RoutingTables3.Attr
				//Delete previous routing data
				database.DeleteRoutingTables(sqliteDatabase,hosts[i].Ip)
			}
							//using previous routingentries if within time

			if len(routingtables) <= 0 {
				log.Print("Routingtables empty")
				continue //routingtables emtpy, try next host
			}

			/*chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,phpsessid)
			if err!= nil {
				chassisType, serialNumber = "database failure", "database failure"
				log.Print(err)
			}*/

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
						log.Print("Error getAPIData, routingentry: ", err)
						continue
					}
					re := &routingEntries{}
					xml.Unmarshal(data2, &re) //Converting XML data to variables
					if err!= nil {
						log.Print("XML error routingentry", err)
						continue
					}
					routingEntries := re.RoutingEntry2.RoutingEntry3.Attr

					entries := regexp.MustCompile(`\d+$`)

					//Because routingentries from the hosts are displayed as a list of for example "2:4", "2:5", we are using regex to get only the routingentries
					for k := range routingEntries {
						tmp := entries.FindStringSubmatch(routingEntries[k])
						for l := range tmp {
							match = append(match, tmp[l])
							//log.Print(tmp[l])
						}
					}
					now := time.Now().Format(time.RFC3339)
					//log.Print("NOW:", now)

					err = database.StoreRoutingEntries(sqliteDatabase, hosts[i].Ip, now, routingtables[j], match)
					if err != nil {
						log.Print(err)
					}
				}

				if (len(match) <= 0) {
						continue
				}
				for k := range match {

					url := "https://" + hosts[i].Ip + "/rest/routingtable/" + routingtables[j] + "/routingentry/" + match[k] + "/historicalstatistics/1"
					_, data3, err := http.GetAPIData(url, phpsessid)
					if err != nil {
						log.Print(err)
						continue
					}

					rData := &rSBCdata{}
					xml.Unmarshal(data3, &rData) //Converting XML data to variables
					if err!= nil {
						log.Print("XML error routing", err)
						continue
					}
					//log.Print("Successful API call data: ",rData.RoutingData)

					metricValue1 := float64(rData.RoutingData.Rt_RuleUsage)
					metricValue2 := float64(rData.RoutingData.Rt_ASR)
					metricValue3 := float64(rData.RoutingData.Rt_RoundTripDelay)
					metricValue4 := float64(rData.RoutingData.Rt_Jitter)
					metricValue5 := float64(rData.RoutingData.Rt_MOS)
					metricValue6 := float64(rData.RoutingData.Rt_QualityFailed)

					c <- prometheus.MustNewConstMetric(collector.Rt_RuleUsage, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname, routingtables[j], match[k])
					c <- prometheus.MustNewConstMetric(collector.Rt_ASR, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname, routingtables[j], match[k])
					c <- prometheus.MustNewConstMetric(collector.Rt_RoundTripDelay, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname, routingtables[j], match[k])
					c <- prometheus.MustNewConstMetric(collector.Rt_Jitter, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname, routingtables[j], match[k])
					c <- prometheus.MustNewConstMetric(collector.Rt_MOS, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname, routingtables[j], match[k])
					c <- prometheus.MustNewConstMetric(collector.Rt_QualityFailed, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname, routingtables[j], match[k])
				}
			}
	}
}

func RoutingEntryCollector() {
	if len(config.GetIncludedHosts("routingentry")) <= 0 {
		log.Print("no hosts")
		return
	}
	c := routingCollector()
	prometheus.MustRegister(c)
}
