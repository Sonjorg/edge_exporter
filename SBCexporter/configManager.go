package main

type hostConfig struct {
	index          int
	ipaddress      string
	systemExporter bool
	Exporter1      bool
	Exporter2      bool
}

import (
    "gopkg.in/yaml.v2"
    )
    //Template used for struct and NewConfig(): https://dev.to/koddr/let-s-write-config-for-your-golang-web-app-on-right-way-yaml-5ggp
    //Specific config is accessed in the format: config.Server.Timeout.Read
    type hostConfig struct {
        //index          int
        HostIp         string `yaml:"host"`
        hostName       string `yaml:"hostname"`
        port           string `yaml:"port"`
        ipaddress      string `yaml:"ip"`
        //exclude        string `yaml:"exclude"`
        Exclude struct {
            // Server is the general server timeout to use
            // for graceful shutdowns
            systemExporter bool `"yaml:systemstats"`
            callStats      bool `"yaml:callstats"`
        }`yaml:"exclude:'true'/'false'"`
    }

    // NewConfig returns a new decoded Config struct
    func NewConfig(configPath string) (*Config, error) {
        // Create config structure
        config := &Config{}
        // Open config file
        file, err := os.Open(configPath)
        if err != nil {
            return nil, err
        }
        defer file.Close()
        // Init new YAML decode
        d := yaml.NewDecoder(file)
        // Start YAML decoding from file
        if err := d.Decode(&config); err != nil {
            return nil, err
        }
        return config, nil
    }
    test := NewConfig(.\config).
    type hosts []hostConfig



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
