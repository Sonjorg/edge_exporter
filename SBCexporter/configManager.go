package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
)

/*
type hostConfig struct {
	index          int
	ipaddress      string
	systemExporter bool
	Exporter1      bool
	Exporter2      bool
}

*/
//Template used for struct and NewConfig(): https://dev.to/koddr/let-s-write-config-for-your-golang-web-app-on-right-way-yaml-5ggp
    type Config struct {
        Hosts []Hosts `yaml:"hosts"`
    }
        //index          int
        type Hosts struct {
            hostName       string `yaml:"hostname"`
            ipaddress         string `yaml:"host"`
            port           string `yaml:"port"`
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
   // test := NewConfig(.\config).
   // type hosts []hostConfig

func getIpAdrExp(exporterName string) []string{
    cfg, err := NewConfig(config.yml)
    if err != nil {
       // log.Fatal(err)
    }
	var list []string
    switch exporterName {
        case "systemStats":
            for i := 0; i < len(cfg.Hosts); i++ {
                if (cfg.Hosts[i].Exclude.systemExporter == false) {
                    list = append(list, cfg.Hosts[i].ipaddress)
                }
            }
        case "callStats":
            for i := 0; i < len(cfg.Hosts); i++ {
                if (cfg.Hosts[i].Exclude.systemExporter == false) {
                    list = append(list, cfg.Hosts[i].ipaddress)
                }
            }
            //INFO: have a switch case on all exporters made, NB!: must remember exact exporternames inside each exporter
        }


return list
}
func main() {
    ip := getIpAdrExp("systemStats")
    fmt.Println(ip)
}
/*func getIPNotExl(exporterName string, hosts *Config) []string {
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
}*/

/*
if hosts[i].Exporter3Excl == true {
	//Exporter2()
	return ipaddr == hosts[i].ipaddress
}
if hosts[i].systemExcl == true {
	return ipaddr == hosts[i].ipaddress
}*/
