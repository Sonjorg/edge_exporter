// routingentry
package collector

import (
	"encoding/xml"
	//"fmt"
	"log"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
)

// /rest/ethernetport
type eSBCdata struct {
	XMLname       xml.Name      `xml:"root"`
	EthernetData  ethernetData  `xml:"ethernetport"`
}
type ethernetData struct {
Href                          string `xml:"href,attr"`
IfName		                  string `xml:"ifName"`
IfAlias                       string `xml:"ifAlias"`
IfRedundancy                  int    `xml:"ifRedundancy"`
IfRedundantPort               int    `xml:"ifRedundantPort"`
Rt_ifInBroadcastPkts		  int    `xml:"rt_ifInBroadcastPkts"`
Rt_ifInDiscards		          int    `xml:"rt_ifInDiscards"`
Rt_ifInErrors		          int    `xml:"rt_ifInErrors"`
Rt_ifInFCSErrors		      int    `xml:"rt_ifInFCSErrors"`
Rt_ifInFragmentedPkts		  int    `xml:"rt_ifInFragmentedPkts"`
Rt_ifInMulticastPkts		  int    `xml:"rt_ifInMulticastPkts"`
Rt_ifInOctets		          int    `xml:"rt_ifInOctets"`
Rt_ifInOverSizedPkts		  int    `xml:"rt_ifInOverSizedPkts"`
Rt_ifInUcastPkts		      int    `xml:"rt_ifInUcastPkts"`
Rt_ifInUndersizedPkts		  int    `xml:"rt_ifInUndersizedPkts"`
Rt_ifInUnknwnProto		      int    `xml:"rt_ifInUnknwnProto"`
Rt_ifInterfaceIndex	          int    `xml:"rt_ifInterfaceIndex"`
Rt_ifLastChange		          int    `xml:"rt_ifLastChange"`
Rt_ifMtu		              int    `xml:"rt_ifMtu"`
Rt_ifName		              int    `xml:"rt_ifName"`
Rt_ifOperatorStatus		      int    `xml:"rt_ifOperatorStatus"`
Rt_ifOutBroadcastPkts		  int    `xml:"rt_ifOutBroadcastPkts"`//Displays the number of transmitted broadcast packets on this port.
Rt_ifOutDeferredTransmissions int    `xml:"rt_ifOutDeferredTransmissions"`//Displays the number of Deferred Transmission errors detected on this port.
Rt_ifOutDiscards		      int    `xml:"rt_ifOutDiscards"` //Displays the number of discard errors detected on this port.
Rt_ifOutErrors		          int    `xml:"rt_ifOutErrors"` //Displays the number of errors detected on this port.
Rt_ifOutLateCollissions		  int    `xml:"rt_ifOutLateCollissions"` //Displays the number of Late Collision errors detected on this port.
Rt_ifOutMulticastPkts		  int    `xml:"rt_ifOutMulticastPkts"` //Displays the number of transmitted multicast packets on this port.
Rt_ifOutOctets		          int    `xml:"rt_ifOutOctets"` //Displays the number of transmitted octets on this port.
Rt_ifOutUcastPkts	    	  int    `xml:"rt_ifOutUcastPkts"` //Displays the number of transmitted unicast packets on this port.
Rt_ifSpeed	                  int    `xml:"rt_ifSpeed"`
Rt_redundancyRole		      int    `xml:"rt_redundancyRole"`
Rt_redundancyState		      int    `xml:"rt_redundancyState"`
}

//Metrics
type ethernetMetrics struct {
IfRedundancy	              *prometheus.Desc
IfRedundantPort	              *prometheus.Desc
Rt_ifInBroadcastPkts	      *prometheus.Desc
Rt_ifInDiscards	              *prometheus.Desc
Rt_ifInErrors	              *prometheus.Desc
Rt_ifInFCSErrors	          *prometheus.Desc
Rt_ifInFragmentedPkts	      *prometheus.Desc
Rt_ifInMulticastPkts	      *prometheus.Desc
Rt_ifInOctets	              *prometheus.Desc
Rt_ifInOverSizedPkts	      *prometheus.Desc
Rt_ifInUcastPkts	          *prometheus.Desc
Rt_ifInUndersizedPkts	      *prometheus.Desc
Rt_ifInUnknwnProto	          *prometheus.Desc
Rt_ifInterfaceIndex	          *prometheus.Desc
Rt_ifLastChange	              *prometheus.Desc
Rt_ifMtu	                  *prometheus.Desc
Rt_ifOperatorStatus	          *prometheus.Desc
Rt_ifOutBroadcastPkts	      *prometheus.Desc
Rt_ifOutDeferredTransmissions *prometheus.Desc
Rt_ifOutDiscards	          *prometheus.Desc
Rt_ifOutErrors	              *prometheus.Desc
Rt_ifOutLateCollissions	      *prometheus.Desc
Rt_ifOutMulticastPkts	      *prometheus.Desc
Rt_ifOutOctets	              *prometheus.Desc
Rt_ifOutUcastPkts	          *prometheus.Desc
Rt_ifSpeed	                  *prometheus.Desc
Rt_redundancyRole 	          *prometheus.Desc
Rt_redundancyState		      *prometheus.Desc
Error_ip                      *prometheus.Desc
}

func ethernetCollector()*ethernetMetrics{

	 return &ethernetMetrics{
		IfRedundancy: prometheus.NewDesc("ifRedundancy",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		IfRedundantPort: prometheus.NewDesc("ifRedundantPort",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInBroadcastPkts: prometheus.NewDesc("rt_ifInBroadcastPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInDiscards: prometheus.NewDesc("rt_ifInDiscards",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInErrors: prometheus.NewDesc("rt_ifInErrors",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInFCSErrors: prometheus.NewDesc("rt_ifInFCSErrors",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInFragmentedPkts: prometheus.NewDesc("rt_ifInFragmentedPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInMulticastPkts: prometheus.NewDesc("rt_ifInMulticastPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInOctets: prometheus.NewDesc("rt_ifInOctets",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInOverSizedPkts: prometheus.NewDesc("rt_ifInOverSizedPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInUcastPkts: prometheus.NewDesc("rt_ifInUcastPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInUndersizedPkts: prometheus.NewDesc("rt_ifInUndersizedPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInUnknwnProto: prometheus.NewDesc("rt_ifInUnknwnProto",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifInterfaceIndex: prometheus.NewDesc("rt_ifInterfaceIndex",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifLastChange: prometheus.NewDesc("rt_ifLastChange",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifMtu: prometheus.NewDesc("rt_ifMtu",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOperatorStatus: prometheus.NewDesc("rt_ifOperatorStatus",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutBroadcastPkts: prometheus.NewDesc("rt_ifOutBroadcastPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutDeferredTransmissions: prometheus.NewDesc("rt_ifOutDeferredTransmissions",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutDiscards: prometheus.NewDesc("rt_ifOutDiscards",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutErrors: prometheus.NewDesc("rt_ifOutErrors",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutLateCollissions: prometheus.NewDesc("rt_ifOutLateCollissions",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutMulticastPkts: prometheus.NewDesc("rt_ifOutMulticastPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutOctets: prometheus.NewDesc("rt_ifOutOctets",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifOutUcastPkts: prometheus.NewDesc("rt_ifOutUcastPkts",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_ifSpeed: prometheus.NewDesc("rt_ifSpeed",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_redundancyRole: prometheus.NewDesc("rt_redundancyRole",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Rt_redundancyState: prometheus.NewDesc("rt_redundancyState",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","ethernetportID","ifName","ifAlias","chassis_type","serial_number"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_ethernetport",
			"NoDescriptionYet",
			[]string{"Instance", "hostname","job","ethernetportID","error_reason","chassis_type","serial_number"}, nil,
		),
	 }
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *ethernetMetrics) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.IfRedundancy
	ch <- collector.IfRedundantPort
	ch <- collector.Rt_ifInBroadcastPkts
	ch <- collector.Rt_ifInDiscards
	ch <- collector.Rt_ifInErrors
	ch <- collector.Rt_ifInFCSErrors
	ch <- collector.Rt_ifInFragmentedPkts
	ch <- collector.Rt_ifInMulticastPkts
	ch <- collector.Rt_ifInOctets
	ch <- collector.Rt_ifInOverSizedPkts
	ch <- collector.Rt_ifInUcastPkts
	ch <- collector.Rt_ifInUndersizedPkts
	ch <- collector.Rt_ifInUnknwnProto
	ch <- collector.Rt_ifInterfaceIndex
	ch <- collector.Rt_ifLastChange
	ch <- collector.Rt_ifMtu
	ch <- collector.Rt_ifOperatorStatus
	ch <- collector.Rt_ifOutBroadcastPkts
	ch <- collector.Rt_ifOutDeferredTransmissions
	ch <- collector.Rt_ifOutDiscards
	ch <- collector.Rt_ifOutErrors
	ch <- collector.Rt_ifOutLateCollissions
	ch <- collector.Rt_ifOutMulticastPkts
	ch <- collector.Rt_ifOutOctets
	ch <- collector.Rt_ifOutUcastPkts
	ch <- collector.Rt_ifSpeed
	ch <- collector.Rt_redundancyRole
	ch <- collector.Rt_redundancyState
	ch <- collector.Error_ip
}
//Collect implements required collect function for all promehteus collectors
func (collector *ethernetMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("ethernetport")//retrieving targets for this exporter
	if (len(hosts) <= 0) {
		log.Print("no hosts")
		return
	}

	for i := range hosts {

		phpsessid,err := http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			log.Print("Error session cookie", hosts[i].Ip, err)
			continue
		}

		//chassis labels from db or http
		chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,phpsessid)
		if err!= nil {
			chassisType, serialNumber = "db chassisData fail", "db chassisData fail"
			log.Print(err)
		}

		var ethernetportID []string
			//Every router has these ethernetports regardless of SBC1000 or SBC2000
			ethernetportID = append(ethernetportID, "23")
			ethernetportID = append(ethernetportID, "29")
			ethernetportID = append(ethernetportID, "24")
			//Couldnt fetch chassis type from db or http: try next host
			for j := range ethernetportID {
					url := "https://"+hosts[i].Ip+"/rest/ethernetport/"+ethernetportID[j]
					_, data, err := http.GetAPIData(url, phpsessid)
						if err != nil {
							log.Print(err)
							c <- prometheus.MustNewConstMetric(
								collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip,hosts[i].Hostname, "ethernetport",ethernetportID[j], "fetching data from this ethernetport failed",chassisType,serialNumber)
							continue
						}

					eData := &eSBCdata{}
					err = xml.Unmarshal(data, &eData) //Converting XML data to variables
					if err!= nil {
						log.Print("XML error ethernet", err)
						continue
					}
					//log.Print("Successful API call data: ",eData.EthernetData)

					metricValue1 := float64(eData.EthernetData.IfRedundancy)
					metricValue2 := float64(eData.EthernetData.IfRedundantPort)
					metricValue3 := float64(eData.EthernetData.Rt_ifInBroadcastPkts)
					metricValue4 := float64(eData.EthernetData.Rt_ifInDiscards)
					metricValue5 := float64(eData.EthernetData.Rt_ifInErrors)
					metricValue6 := float64(eData.EthernetData.Rt_ifInFCSErrors)
					metricValue7 := float64(eData.EthernetData.Rt_ifInFragmentedPkts)
					metricValue8 := float64(eData.EthernetData.Rt_ifInMulticastPkts)
					metricValue9 := float64(eData.EthernetData.Rt_ifInOctets)
					metricValue10 := float64(eData.EthernetData.Rt_ifInOverSizedPkts)
					metricValue11 := float64(eData.EthernetData.Rt_ifInUcastPkts)
					metricValue12 := float64(eData.EthernetData.Rt_ifInUndersizedPkts)
					metricValue13 := float64(eData.EthernetData.Rt_ifInUnknwnProto)
					metricValue14 := float64(eData.EthernetData.Rt_ifInterfaceIndex)
					metricValue15 := float64(eData.EthernetData.Rt_ifLastChange)
					metricValue16 := float64(eData.EthernetData.Rt_ifMtu)
					metricValue18 := float64(eData.EthernetData.Rt_ifOperatorStatus)
					metricValue19 := float64(eData.EthernetData.Rt_ifOutBroadcastPkts)
					metricValue20 := float64(eData.EthernetData.Rt_ifOutDeferredTransmissions)
					metricValue21 := float64(eData.EthernetData.Rt_ifOutDiscards)
					metricValue22 := float64(eData.EthernetData.Rt_ifOutErrors)
					metricValue23 := float64(eData.EthernetData.Rt_ifOutLateCollissions)
					metricValue24 := float64(eData.EthernetData.Rt_ifOutMulticastPkts)
					metricValue25 := float64(eData.EthernetData.Rt_ifOutOctets)
					metricValue26 := float64(eData.EthernetData.Rt_ifOutUcastPkts)
					metricValue27 := float64(eData.EthernetData.Rt_ifSpeed)
					metricValue28 := float64(eData.EthernetData.Rt_redundancyRole)
					metricValue29 := float64(eData.EthernetData.Rt_redundancyState)
						c <- prometheus.MustNewConstMetric(collector.IfRedundancy, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.IfRedundantPort, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInBroadcastPkts, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInDiscards, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInErrors, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInFCSErrors, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInFragmentedPkts, prometheus.GaugeValue, metricValue7, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInMulticastPkts, prometheus.GaugeValue, metricValue8, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInOctets, prometheus.GaugeValue, metricValue9, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInOverSizedPkts, prometheus.GaugeValue, metricValue10, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInUcastPkts, prometheus.GaugeValue, metricValue11, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInUndersizedPkts, prometheus.GaugeValue, metricValue12, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInUnknwnProto, prometheus.GaugeValue, metricValue13, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifInterfaceIndex, prometheus.GaugeValue, metricValue14, hosts[i].Ip, hosts[i].Hostname,"ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifLastChange, prometheus.GaugeValue, metricValue15, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifMtu, prometheus.GaugeValue, metricValue16, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOperatorStatus, prometheus.GaugeValue, metricValue18, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutBroadcastPkts, prometheus.GaugeValue, metricValue19, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutDeferredTransmissions, prometheus.GaugeValue, metricValue20, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutDiscards, prometheus.GaugeValue, metricValue21, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutErrors, prometheus.GaugeValue, metricValue22, hosts[i].Ip, hosts[i].Hostname,"ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutLateCollissions, prometheus.GaugeValue, metricValue23, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutMulticastPkts, prometheus.GaugeValue, metricValue24, hosts[i].Ip, hosts[i].Hostname,"ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutOctets, prometheus.GaugeValue, metricValue25, hosts[i].Ip, hosts[i].Hostname,"ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifOutUcastPkts, prometheus.GaugeValue, metricValue26, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ifSpeed, prometheus.GaugeValue, metricValue27, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_redundancyRole, prometheus.GaugeValue, metricValue28, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_redundancyState, prometheus.GaugeValue, metricValue29, hosts[i].Ip, hosts[i].Hostname, "ethernetport",ethernetportID[j],eData.EthernetData.IfName,eData.EthernetData.IfAlias, chassisType, serialNumber)
		}
	}
}

// Initializing the collector
func EthernetportCollector() {
	hosts := config.GetIncludedHosts("ethernetport")//retrieving targets for this exporter
	if (len(hosts) <= 0) {
		//log.Print("no hosts ethernetport")
		return
	}
		c := ethernetCollector()
		prometheus.MustRegister(c)
}
