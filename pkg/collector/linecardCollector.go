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

// /rest/linecard
type lSBCdata struct {
	XMLname       xml.Name      `xml:"root"`
	LinecardData  LinecardData  `xml:"linecard"`
}
type LinecardData struct {
Href              string `xml:"href,attr"`
Rt_CardType		  int    `xml:"rt_CardType"`
Rt_Location		  int    `xml:"rt_Location"`
Rt_ServiceStatus  int    `xml:"rt_ServiceStatus"`
Rt_Status         int    `xml:"rt_Status"`

}

//Metrics for each routingentry
type linecardMetrics struct {
	Href                *prometheus.Desc
	Rt_CardType		  	*prometheus.Desc
	Rt_Location		  	*prometheus.Desc
	Rt_ServiceStatus  	*prometheus.Desc
	Rt_Status           *prometheus.Desc
	Error_ip            *prometheus.Desc
	}

func lineCCollector()*linecardMetrics{

	 return &linecardMetrics{
		Rt_CardType: prometheus.NewDesc("rt_CardType",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","linecardID","chassis_type","serial_number"}, nil,
		),
		Rt_Location: prometheus.NewDesc("rt_Location",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","linecardID","chassis_type","serial_number"}, nil,
		),
		Rt_ServiceStatus: prometheus.NewDesc("rt_ServiceStatus",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","linecardID","chassis_type","serial_number"}, nil,
		),
		Rt_Status: prometheus.NewDesc("rt_Status",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","linecardID","chassis_type","serial_number"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_linecard",
			"NoDescriptionYet",
			[]string{"Instance", "hostname","job","linecardID","error_reason","chassis_type","serial_number"}, nil,
		),
	 }
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *linecardMetrics) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.Rt_CardType
	ch <- collector.Rt_Location
	ch <- collector.Rt_ServiceStatus
	ch <- collector.Rt_Status
	ch <- collector.Error_ip
}
//Collect implements required collect function for all promehteus collectors
func (collector *linecardMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("linecard")//retrieving targets for this exporter
	if (len(hosts) <= 0) {
		log.Print("no hosts")
		return
	}
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64
	var metricValue4 float64

	for i := range hosts {

		phpsessid,err := http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			log.Print("Error auth", hosts[i].Ip, err)
			continue
		}

		//chassis labels from db or http
		chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,phpsessid)
		if err!= nil {
			chassisType, serialNumber = "db chassisData fail", "db chassisData fail"
			log.Print(err)
		}

		var linecardID []string
		// There are two linecard linecardIDs which are different for type of SBC router
		if (chassisType == "SBC1000") {
			linecardID = append(linecardID, "7")
			linecardID = append(linecardID, "8")
		} else if (chassisType == "SBC2000") {
			linecardID = append(linecardID, "1")
			linecardID = append(linecardID, "2")
		} else {
			//Couldnt fetch chassis type from db or http: try next host
			c <- prometheus.MustNewConstMetric(
				collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip,hosts[i].Hostname, "linecard","All linecards", "fetching chassistype from db or http fail","NA","NA")
			continue
		}
			for j := range linecardID {
					url := "https://"+hosts[i].Ip+"/rest/linecard/"+linecardID[j]
					_, data, err := http.GetAPIData(url, phpsessid)
						if err != nil {
							log.Print(err)
							c <- prometheus.MustNewConstMetric(
								collector.Error_ip, prometheus.GaugeValue, 0, hosts[i].Ip,hosts[i].Hostname, "linecard",linecardID[j], "fetching data from this linecard failed",chassisType,serialNumber)
							continue
						}

					lData := &lSBCdata{}
					err = xml.Unmarshal(data, &lData) //Converting XML data to variables
					//log.Print("Successful API call data: ",lData.LinecardData,"\n")
					if err!= nil {
						log.Print("XML error linecard", err)
						//continue
					}
					metricValue1 = float64(lData.LinecardData.Rt_CardType)
					metricValue2 = float64(lData.LinecardData.Rt_Location)
					metricValue3 = float64(lData.LinecardData.Rt_ServiceStatus)
					metricValue4 = float64(lData.LinecardData.Rt_Status)

						c <- prometheus.MustNewConstMetric(collector.Rt_CardType, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname, "linecard",linecardID[j], chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_Location, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname, "linecard",linecardID[j], chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_ServiceStatus, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname, "linecard",linecardID[j], chassisType, serialNumber)
						c <- prometheus.MustNewConstMetric(collector.Rt_Status, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname, "linecard",linecardID[j], chassisType, serialNumber)


		}
	}
}

// Initializing the collector
func LinecardCollector() {
		c := lineCCollector()
		prometheus.MustRegister(c)
}
