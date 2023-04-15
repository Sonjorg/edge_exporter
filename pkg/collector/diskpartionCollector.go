// routingentry
package collector

import (
	"encoding/xml"
	"fmt"

	//"log"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
)


 type diskPartition struct {
	XMLName           xml.Name          `xml:"root"`
	DiskPartitionList diskPartitionList `xml:"diskpartition_list"`
 }
 type diskPartitionList struct {
	DiskPartitionEntry diskPartitionEntry  `xml:"diskpartition_pk"`
 }
 type diskPartitionEntry struct {
	Attr    []string `xml:"id,attr"`
	Value     string `xml:",chardata"`
 }

//second request
type dSBCdata struct {
	XMLname   xml.Name  `xml:"root"`
	DiskData  diskData  `xml:"diskpartition"`
}
type diskData struct {
Href                string `xml:"href,attr"`
Rt_CurrentUsage		int    `xml:"rt_CurrentUsage"`
Rt_MaximumSize		int    `xml:"rt_MaximumSize"`
Rt_MemoryAvailable	int    `xml:"rt_MemoryAvailable"`
Rt_MemoryUsed       int    `xml:"rt_MemoryUsed"`
Rt_PartitionName    string `xml:"rt_PartitionName"`
Rt_PartitionType    int    `xml:"rt_PartitionType"`
}

//Metrics for each routingentry
type diskMetrics struct {
	Href                *prometheus.Desc
	Rt_CurrentUsage		*prometheus.Desc
	Rt_MaximumSize		*prometheus.Desc
	Rt_MemoryAvailable	*prometheus.Desc
	Rt_MemoryUsed       *prometheus.Desc
	//Rt_PartitionName    *prometheus.Desc
	Rt_PartitionType    *prometheus.Desc
	Error_ip            *prometheus.Desc
	}

func diskCollector()*diskMetrics{

	 return &diskMetrics{
		Rt_CurrentUsage: prometheus.NewDesc("rt_CurrentUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","disk_partition_id","disk_partition_name","chassis_type","serial_number"}, nil,
		),
		Rt_MaximumSize: prometheus.NewDesc("rt_MaximumSize",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","disk_partition_id","disk_partition_name","chassis_type","serial_number"}, nil,
		),
		Rt_MemoryAvailable: prometheus.NewDesc("rt_MemoryAvailable",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","disk_partition_id","disk_partition_name","chassis_type","serial_number"}, nil,
		),
		Rt_MemoryUsed: prometheus.NewDesc("rt_MemoryUsed",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","disk_partition_id","disk_partition_name","chassis_type","serial_number"}, nil,
		),
		Rt_PartitionType: prometheus.NewDesc("rt_PartitionType",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job","disk_partition_id","disk_partition_name","chassis_type","serial_number"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_edge_disk",
			"NoDescriptionYet",
			[]string{"Instance", "hostname"}, nil,
		),
	 }
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *diskMetrics) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.Rt_CurrentUsage
	ch <- collector.Rt_MaximumSize
	ch <- collector.Rt_MemoryAvailable
	ch <- collector.Rt_MemoryUsed
	//ch <- collector.Rt_PartitionName
	ch <- collector.Rt_PartitionType
	ch <- collector.Error_ip
}
//Collect implements required collect function for all promehteus collectors
func (collector *diskMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("diskpartition")//retrieving targets for this exporter
	if (len(hosts) <= 0) {
		fmt.Println("no hosts")
		return
	}
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64
	var metricValue4 float64
	var metricValue5 float64
	var partitionName string

	for i := range hosts {

		phpsessid,err := http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			fmt.Println("Error auth", hosts[i].Ip, err)
			continue
		}
		_, data,err := http.GetAPIData("https://"+hosts[i].Ip+"/rest/diskpartition", phpsessid)
		if err != nil {
			fmt.Println("Error disk data", hosts[i].Ip, err)
			continue
		}
		//b := []byte(data) //Converting string of data to bytestream
		disk := &diskPartition{}
		xml.Unmarshal(data, &disk) //Converting XML data to variables

		//List of disks retrieved from routers as XML
		disks := disk.DiskPartitionList.DiskPartitionEntry.Attr
		if (len(disks) <= 0) {
			//return nil, "Routingtables empty"
			fmt.Println("disks empty")
			continue

		}
		//chassis labels from db or http
		chassisType, serialNumber, err := utils.GetChassisLabels(hosts[i].Ip,phpsessid)
		if err!= nil {
			chassisType, serialNumber = "db chassisData fail", "db chassisData fail"
			fmt.Println(err)
		}
			for j := range disks {

					url := "https://"+hosts[i].Ip+"/rest/diskpartition/"+disks[j]
					_, data2, err := http.GetAPIData(url, phpsessid)
						if err != nil {
							fmt.Println(err)

							continue
						}

					dData := &dSBCdata{}
					err = xml.Unmarshal(data2, &dData) //Converting XML data to variables
					fmt.Println("Successful API call data: ",dData.DiskData)
					if err!= nil {
						fmt.Println("XML error disk", err)
						//continue
					}
					metricValue1 = float64(dData.DiskData.Rt_CurrentUsage)
					metricValue2 = float64(dData.DiskData.Rt_MaximumSize)
					metricValue3 = float64(dData.DiskData.Rt_MemoryAvailable)
					metricValue4 = float64(dData.DiskData.Rt_MemoryUsed)
					metricValue5 = float64(dData.DiskData.Rt_PartitionType)
					partitionName = dData.DiskData.Rt_PartitionName
					id := string(j)

						c <- prometheus.MustNewConstMetric(collector.Rt_CurrentUsage, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname, "diskpartition",id, partitionName,chassisType, serialNumber,)
						c <- prometheus.MustNewConstMetric(collector.Rt_MaximumSize, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname, "diskpartition",id, partitionName,chassisType, serialNumber,)
						c <- prometheus.MustNewConstMetric(collector.Rt_MemoryAvailable, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname, "diskpartition",id, partitionName,chassisType, serialNumber,)
						c <- prometheus.MustNewConstMetric(collector.Rt_MemoryUsed, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname, "diskpartition",id, partitionName,chassisType, serialNumber,)
						c <- prometheus.MustNewConstMetric(collector.Rt_PartitionType, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname, "diskpartition",id, partitionName,chassisType, serialNumber,)


		}
	}
}

// Initializing the exporter


func DiskPartitionCollector() {
		c := diskCollector()
		prometheus.MustRegister(c)
}
