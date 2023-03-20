package main

import (
	"fmt"
	//"os"
	"gopkg.in/yaml.v2"
    //"flag"
   // "log"
   "io/ioutil"
)
// Template used for struct and the functions NewConfig(), ValidateConfigPath() and ParseFlags() are copied from:
// https://dev.to/koddr/let-s-write-config-for-your-golang-web-app-on-right-way-yaml-5ggp

    type Config struct {
        Hosts []Host
        Authtimeout int `yaml:"authtimeout"`
    }
    type Host struct {
        HostName       string `yaml:"hostname"`
        Ipaddress      string `yaml:"ipaddress"`
        Username       string `yaml:"username"`
        Password       string `yaml:"password"`
        Exclude      []string `yaml:"exclude"`
        }

            //From stackoverflow
    func getConf(c *Config) *Config {
        yamlFile, err := ioutil.ReadFile("../config.yml")
            if err != nil {
                  //log.Printf("yamlFile.Get err   #%v ", err)
                     fmt.Println("yamlFile.Get err   # ", err)
            }
        err = yaml.Unmarshal(yamlFile, c)
        if err != nil {
             // log.Fatalf("Unmarshal: %v", err)
              fmt.Println("yamlFile.Get err   # ", err)
         }
      return c
     }

     type includedHosts struct {
        ip         string
        hostname   string
        username   string
        password   string
    }
    // This functions iterates through all hosts in the saved config and
    // returns a list of hosts that doesn't have the specified collector excluded in the config file
    // exporterName must be equal to "system", "routingentry" ..
     func getIncludedHosts(collectorName string) []includedHosts {
        cfg := getConf(&Config{})
        list := make([]includedHosts,0,8)
        var excluded bool

        for i := range cfg.Hosts {
            for v := range cfg.Hosts[i].Exclude {
                if (cfg.Hosts[i].Exclude[v] == collectorName) {
                    excluded = true
                }
            }
            if !excluded {
                list = append(list, includedHosts{cfg.Hosts[i].Ipaddress, cfg.Hosts[i].HostName,cfg.Hosts[i].Username, cfg.Hosts[i].Password})
            }
        }
    return list
    }

/*

func main() {

    g:= getIncludedHosts("system")
    for i:= range g {
    fmt.Println(g[i].hostname,g[i].username)
    }
}
*/