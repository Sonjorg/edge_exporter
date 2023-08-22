/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
*/
package collector

import (
	"encoding/xml"
	"log"
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
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

func DiskPartitionCollector(host *config.HostCompose)(m []prometheus.Metric) {

	var (
		Rt_CurrentUsage = prometheus.NewDesc("edge_disk_CurrentUsage",
			"Amount of memory used by this partition, expressed as percentage",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_MaximumSize = prometheus.NewDesc("edge_disk_MaximumSize",
			"Specifies the maximum amount of memory, in bytes available in this partition.",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_MemoryAvailable = prometheus.NewDesc("edge_disk_Available",
			"Amount of memory in bytes, available for use in the filesystem.",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_MemoryUsed = prometheus.NewDesc("edge_disk_Used",
			"Amount of memory in bytes, used by the existing files in the filesystem",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
		Rt_PartitionType = prometheus.NewDesc("edge_disk_PartitionType",
			"Identifies the user-friendly physical device holding the partition.",
			[]string{"hostip", "hostname", "disk_partition_id","disk_partition_name"}, nil,
		)
	)


		phpsessid,err := http.APISessionAuth(host.Username, host.Password, host.Ip)
		if err != nil {
			log.Print("Error auth", host.Ip, err)
			return
		}

		_, data,err := http.GetAPIData("https://"+host.Ip+"/rest/diskpartition", phpsessid)
		if err != nil {
			log.Print("Error fetching diskpartition data: ", host.Ip, err)
			return
		}
		disk := &diskPartition{}
		xml.Unmarshal(data, &disk) //Converting XML data to variables

		//List of disks retrieved
		disks := disk.DiskPartitionList.DiskPartitionEntry.Attr
		if (len(disks) <= 0) {
			log.Print("disks empty")
			return

		}

			for j := range disks {

					url := "https://"+host.Ip+"/rest/diskpartition/"+disks[j]
					_, data2, err := http.GetAPIData(url, phpsessid)
						if err != nil {
							log.Print("Error fetching diskpartition data = ", host.Ip, err)
							continue
						}

					dData := &dSBCdata{}
					err = xml.Unmarshal(data2, &dData) //Converting XML data to variables
					//log.Print("Successful API call data: ",dData.DiskData)
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

					m = append(m, prometheus.MustNewConstMetric(Rt_CurrentUsage, prometheus.GaugeValue, metricValue1, host.Ip, host.Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_MaximumSize, prometheus.GaugeValue, metricValue2, host.Ip, host.Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_MemoryAvailable, prometheus.GaugeValue, metricValue3, host.Ip, host.Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_MemoryUsed, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname,id, partitionName))
					m = append(m, prometheus.MustNewConstMetric(Rt_PartitionType, prometheus.GaugeValue, metricValue5, host.Ip, host.Hostname,id, partitionName))
		}
	return m
}


