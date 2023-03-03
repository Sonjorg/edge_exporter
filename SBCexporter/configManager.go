package main

type hostConfig struct {
	index          int
	ipaddress      string
	systemExporter bool
	Exporter1      bool
	Exporter2      bool
}

type hosts []hostConfig

func getIPNotExl(exporterName string, hosts hostConfig) []string{

	var list []string
	//var i = hosts.index.size
	if (exporterName == "systemExporter") {

		for i := 0; i < hosts.index; i++ {
			if (hosts.systemExporter == true) {
				list = append(list, hosts.ipaddress)
				}
		}
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
