/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
 */
package collector

import (
	"database/sql"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
)

// first request
// rest/routingtable/
type routingTables struct {
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
type routingTableX struct {
	XMLName      xml.Name       `xml:"root"`
	Routingtable routingTableX2 `xml:"routingtable"`
}
type routingTableX2 struct {
	Description string `xml:"Description"`
	Sequence    string `xml:"Sequence"`
}

//Will soon remove the following
// rest/routingtable/4/routingentry

// Third request
// rest/routingtable/2/routingentry/1/historicalstatistics/1
type rSBCdata struct {
	XMLname     xml.Name    `xml:"root"`
	Status      rStatus     `xml:"status"`
	RoutingData routingData `xml:"routingentry"`
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
	Description       string `xml:"Description"`
}

func RoutingEntryCollector(host *config.HostCompose) (m []prometheus.Metric) {

	var (
		Rt_RuleUsage = prometheus.NewDesc("edge_routingentry_RuleUsage",
			"Displays the number of times this call route has been selected for a call.",
			[]string{"hostip", "hostname", "rt_description", "re_description"}, nil, //, "routing_table", "routing_entry"
		)
		Rt_ASR = prometheus.NewDesc("edge_routingentry_ASR",
			"Displays the Answer-Seizure Ratio for this call route. (ASR is calculated by dividing the number of call attempts answered by the number of call attempts.)",
			[]string{"hostip", "hostname", "rt_description", "re_description"}, nil,
		)
		Rt_RoundTripDelay = prometheus.NewDesc("edge_routingentry_RoundTripDelay",
			"Displays the average round trip delay for this call route.",
			[]string{"hostip", "hostname", "rt_description", "re_description"}, nil,
		)
		Rt_Jitter = prometheus.NewDesc("edge_routingentry_Jitter",
			"Displays the average jitter for this call route.",
			[]string{"hostip", "hostname", "rt_description", "re_description"}, nil,
		)
		Rt_MOS = prometheus.NewDesc("edge_routingentry_MOS",
			"Displays the Mean Opinion Score (MOS) for this call route.",
			[]string{"hostip", "hostname", "rt_description", "re_description"}, nil,
		)
		Rt_QualityFailed = prometheus.NewDesc("edge_routingentry_QualityFailed",
			"Displays if this call route is currently passing or failing the associated quality metrics. If true then the rule is failing, if false then it is passing.",
			[]string{"hostip", "hostname", "rt_description", "re_description"}, nil,
		)
	)

	var sqliteDatabase *sql.DB
	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Print(err)
	}
	var rtdescription string
	var timeLastString string //Fetched from database, the routingentries and tables are stored for 24 hours as requested by HDO
	var timeLast time.Time
	phpsessid, err := http.APISessionAuth(host.Username, host.Password, host.Ip)
	if err != nil {
		log.Print("Error authentication", host.Ip, err)
		return
	}
	var routingtables []string
	var routingEntryMap = make(map[string][]string)
	var DBexists bool = database.RoutingTablesExists(sqliteDatabase, host.Ip) //Previous data is stored in db? Fetch this data
	if DBexists {
		routingEntryMap, routingtables, timeLastString, rtdescription, err = database.GetRoutingData(sqliteDatabase, host.Ip) // From db = returning a map of routingentables to routingentries (array),
		if err != nil {
			fmt.Print(err)
		}
		timeLast, err = time.Parse(time.RFC3339, timeLastString)
		if err != nil {
			fmt.Print(err)
		}
	}

	timeSchedule := host.RoutingEntryTime
	//If 24 hours has not passed since last data was stored in database, use this data
	if !DBexists || utils.Expired(timeSchedule, timeLast) { //Routing data has expired, fetching new routingentries
		fmt.Println("Fetching routing data from http")
		_, data, err := http.GetAPIData("https://"+host.Ip+"/rest/routingtable", phpsessid)
		if err != nil {
			fmt.Print("Error routingtable data", host.Ip, err)
			return
		}

		rt := &routingTables{}
		err = xml.Unmarshal(data, &rt) //Converting XML data to variables
		if err != nil {
			fmt.Print("XML error routingentry", err)
			return
		}

		routingtables = rt.RoutingTables2.RoutingTables3.Attr
		//Delete previous routing data
		database.DeleteRoutingTables(sqliteDatabase, host.Ip)
	}

	if len(routingtables) <= 0 {
		fmt.Print("Routingtables empty")
		return //routingtables emtpy, try next host
	}

	for j := range routingtables {
		var routingEntries []string //variable to hold routingentries cleaned with regex

		//Trying to fetch routingentries from database, if not exist yet, fetch new ones
		if DBexists {
			for k, v := range routingEntryMap { // fetching routingentries from a map from db
				if k == routingtables[j] {
					for re := range v {
						routingEntries = append(routingEntries, v[re]) //using previous routingentries (match)
					}
				}
			}
			fmt.Println("Fetching routingdata from database")

		} else { // DB doesn't exist, so fetch new routingentries with
			url := "https://" + host.Ip + "/rest/routingtable/" + routingtables[j]
			_, data2, err := http.GetAPIData(url, phpsessid)
			if err != nil {
				fmt.Print("Error getAPIData, routingentry", err)
				continue
			}
			re := &routingTableX{}
			xml.Unmarshal(data2, &re) //Converting XML data to variables
			if err != nil {
				fmt.Print("XML error routingentry", err)
				continue
			}
			routingE := re.Routingtable.Sequence
			routingEntries = strings.Split(routingE, ",")

			rtdescription = re.Routingtable.Description

			now := time.Now().Format(time.RFC3339)

			err = database.StoreRoutingEntries(sqliteDatabase, host.Ip, now, routingtables[j], routingEntries, rtdescription)
			if err != nil {
				fmt.Print(err)
			}
		}
		if len(routingEntries) <= 0 {
			continue
		}
		for k := range routingEntries {

			url := "https://" + host.Ip + "/rest/routingtable/" + routingtables[j] + "/routingentry/" + routingEntries[k]
			_, data3, err := http.GetAPIData(url, phpsessid)
			if err != nil {
				fmt.Print(err)
				continue
			}

			rData := &rSBCdata{}
			xml.Unmarshal(data3, &rData) //Converting XML data to variables
			if err != nil {
				fmt.Print("XML error routing", err)
				continue
			}

			metricValue1 := float64(rData.RoutingData.Rt_RuleUsage)
			metricValue2 := float64(rData.RoutingData.Rt_ASR)
			metricValue3 := float64(rData.RoutingData.Rt_RoundTripDelay)
			metricValue4 := float64(rData.RoutingData.Rt_Jitter)
			metricValue5 := float64(rData.RoutingData.Rt_MOS)
			metricValue6 := float64(rData.RoutingData.Rt_QualityFailed)
			redesc := rData.RoutingData.Description

			m = append(m, prometheus.MustNewConstMetric(Rt_RuleUsage, prometheus.GaugeValue, metricValue1, host.Ip, host.Hostname, rtdescription, redesc))
			m = append(m, prometheus.MustNewConstMetric(Rt_ASR, prometheus.GaugeValue, metricValue2, host.Ip, host.Hostname, rtdescription, redesc))
			m = append(m, prometheus.MustNewConstMetric(Rt_RoundTripDelay, prometheus.GaugeValue, metricValue3, host.Ip, host.Hostname, rtdescription, redesc))
			m = append(m, prometheus.MustNewConstMetric(Rt_Jitter, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname, rtdescription, redesc))
			m = append(m, prometheus.MustNewConstMetric(Rt_MOS, prometheus.GaugeValue, metricValue5, host.Ip, host.Hostname, rtdescription, redesc))
			m = append(m, prometheus.MustNewConstMetric(Rt_QualityFailed, prometheus.GaugeValue, metricValue6, host.Ip, host.Hostname, rtdescription, redesc)) //, routingtables[j], routingEntries[k]
		}
	}

	return m
}
