// routingentry
package collector

import (
	"encoding/xml"
	//"fmt"
	"log"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	//"edge_exporter/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
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

func DiskPartitionCollector()(m []prometheus.Metric) {
	
	hosts := config.GetIncludedHosts("diskpartition")//retrieving targets for this exporter
	if (len(hosts) <= 0) {
		log.Print("no hosts disk")
		return
	}
	var (
		Rt_CurrentUsage = prometheus.NewDesc("rt_CurrentUsage",
			"Amount of memory used by this partition, expressed as percentage",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_MaximumSize = prometheus.NewDesc("rt_MaximumSize",
			"Specifies the maximum amount of memory, in bytes available in this partition.",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_MemoryAvailable = prometheus.NewDesc("rt_MemoryAvailable",
			"Amount of memory in bytes, available for use in the filesystem.",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_MemoryUsed = prometheus.NewDesc("rt_MemoryUsed",
			"Amount of memory in bytes, used by the existing files in the filesystem",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_PartitionType = prometheus.NewDesc("rt_PartitionType",
			"Identifies the user-friendly physical device holding the partition.",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
	)

	for i := range hosts {

		phpsessid,err := http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
		if err != nil {
			log.Print("Error auth", hosts[i].Ip, err)
			continue
		}

		_, data,err := http.GetAPIData("https://"+hosts[i].Ip+"/rest/diskpartition", phpsessid)
		if err != nil {
			log.Print("Error fetching diskpartition data = ", hosts[i].Ip, err)
			continue
		}
		disk := &diskPartition{}
		xml.Unmarshal(data, &disk) //Converting XML data to variables

		//List of disks retrieved from routers as XML
		disks := disk.DiskPartitionList.DiskPartitionEntry.Attr
		if (len(disks) <= 0) {
			//return nil, "Routingtables empty"
			log.Print("disks empty")
			continue

		}

			for j := range disks {

					url := "https://"+hosts[i].Ip+"/rest/diskpartition/"+disks[j]
					_, data2, err := http.GetAPIData(url, phpsessid)
						if err != nil {
							log.Print("Error fetching diskpartition data = ", hosts[i].Ip, err)
							continue
						}

					dData := &dSBCdata{}
					err = xml.Unmarshal(data2, &dData) //Converting XML data to variables
					//log.Print("Successful API call data = ",dData.DiskData)
					if err!= nil {
						log.Print("XML error disk", err)
						continue
					}
					metricValue1 := float64(dData.DiskData.Rt_CurrentUsage)
					metricValue2 := float64(dData.DiskData.Rt_MaximumSize)
					metricValue3 := float64(dData.DiskData.Rt_MemoryAvailable)
					metricValue4 := float64(dData.DiskData.Rt_MemoryUsed)
					metricValue5 := float64(dData.DiskData.Rt_PartitionType)
					partitionName := dData.DiskData.Rt_PartitionName
					id := strconv.Itoa(j)

					m = append(m, prometheus.MustNewConstMetric(Rt_CurrentUsage, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_MaximumSize, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_MemoryAvailable, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_MemoryUsed, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_PartitionType, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname,id, partitionName))
		}
	}
	return m
}

// Initializing the exporter


