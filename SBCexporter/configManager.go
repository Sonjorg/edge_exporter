package main

type hostConfig struct {
	index          int
	ipaddress      string
	systemExporter bool
	Exporter1      bool
	Exporter2      bool
}

type hosts []hostConfig

func getIpAdrExp(exporterName string, hosts hostConfig) []string{

	var list []string

	for i := 0; i < hosts.index; i++ {
		if (hosts.systemExporter == true) {
			list = append(list, hosts.ipaddress)
			}
	}
return list
}
func getIPNotExl(exporterName string, hosts hostConfig) []string {
	var list []string

	switch exporterName {
	case "systemStatsExporter":
			list = getIpAdrExp(exporterName, hosts)
	//var i = hosts.index.size
	case "teleStatsExporter":
			list = getIpAdrExp(exporterName, hosts)
		//INFO: have a switch case on all exporters made, NB!: must remember exact exporternames inside each exporter
	}
	return list
}

/*
if hosts[i].Exporter3Excl == true {
	//Exporter2()
	return ipaddr == hosts[i].ipaddress
}
if hosts[i].systemExcl == true {
	return ipaddr == hosts[i].ipaddress
}*/
