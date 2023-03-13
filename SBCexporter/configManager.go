package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
    "flag"
)

//Template used for struct and NewConfig(): https://dev.to/koddr/let-s-write-config-for-your-golang-web-app-on-right-way-yaml-5ggp
    type Config struct {
        Hosts []Host
    }
        type Host struct {
            HostName       string `yaml:"hostname"`
            Ipaddress      string `yaml:"ipaddress"`
            //exclude        string `yaml:"exclude"`
                Exclude struct {
                    // Server is the general server timeout to use
                    // for graceful shutdowns
                    SystemExporter bool `yaml:"systemstats"`
                    CallStats      bool `yaml:"callstats"`
                }`yaml:"exclude"`
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

   func ValidateConfigPath(path string) error {
    s, err := os.Stat(path)
    if err != nil {
        return err
    }
    if s.IsDir() {
        return fmt.Errorf("'%s' is a directory, not a normal file", path)
    }
    return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
    // String that contains the configured configuration path
    var configPath string

    // Set up a CLI flag called "-config" to allow users
    // to supply the configuration file
    flag.StringVar(&configPath, "config", "./config.yml", "./config.yml")

    // Actually parse the flags
    flag.Parse()

    // Validate the path first
    if err := ValidateConfigPath(configPath); err != nil {
        return "", err
    }

    // Return the configuration path
    return configPath, nil
}
func getIpAdrExp(exporterName string) []string{
    cfgPath, err := ParseFlags()
    if err != nil {
        fmt.Println(err)
    }
    cfg, err := NewConfig(cfgPath)
    if err != nil {
       fmt.Println(err)
    }
	var list []string
    switch exporterName {
        case "systemStats":
           for i := range cfg.Hosts {
            //for i := 0; i < len(cfg.Hosts); i++ {
                if (cfg.Hosts[i].Exclude.SystemExporter == false) {
                    list = append(list, cfg.Hosts[i].Ipaddress)
                }
            }
        case "callStats":
            for i:= range cfg.Hosts {
                if (cfg.Hosts[i].Exclude.CallStats == false) {
                    list = append(list, cfg.Hosts[i].Ipaddress)
                }
            }
            //INFO: have a switch case on all exporters made, NB!: must remember exact exporternames inside each exporter
        }

return list
}
func test() {
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
