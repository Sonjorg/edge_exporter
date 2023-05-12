package config

import (
	//"fmt"
	//"os"
	"gopkg.in/yaml.v2"
    //"flag"
    "log"
   "io/ioutil"
)

//Describing config.yml file
    type Config struct {
        Hosts []Host
        Authtimeout int `yaml:"authtimeout"`
    }
    type Host struct {
        HostName          string `yaml:"hostname"`
        Ipaddress         string `yaml:"ipaddress"`
        Username          string `yaml:"username"`
        Password          string `yaml:"password"`
        Exclude         []string `yaml:"exclude"`
        RoutingEntryTime  float64 `yaml:"routing-database-hours"`
    }

    //GetConf is from stackoverflow
    func GetConf(c *Config) *Config {
        yamlFile, err := ioutil.ReadFile("config.yml")
            if err != nil {
                     log.Print("yamlFile.Get err   # ", err)
            }
        err = yaml.Unmarshal(yamlFile, c)
        if err != nil {
              log.Print("yamlFile.Get err   # ", err)
         }
      return c
     }

     type IncludedHosts struct {
        Ip         string
        Hostname   string
        Username   string
        Password   string
        RoutingEntryTime  float64
    }

    func GetAllHosts() []IncludedHosts{
        cfg := GetConf(&Config{})
        list := make([]IncludedHosts,0,8)

        for i := range cfg.Hosts {
            list = append(list, IncludedHosts{cfg.Hosts[i].Ipaddress, cfg.Hosts[i].HostName,cfg.Hosts[i].Username, cfg.Hosts[i].Password,cfg.Hosts[i].RoutingEntryTime})
        }
    return list
    }

    // Returns hosts' config for hosts with a collectorname that have not been excluded in config.yml
    // This functions iterates through all hosts in the saved config and
    // returns a list of hosts that doesn't have the specified collector excluded in the config file
    // exporterName must be equal to "system", "routingentry" ..
    func GetIncludedHosts(collectorName string) []IncludedHosts {
        cfg := GetConf(&Config{})
        list := make([]IncludedHosts,0,8)
        var excluded bool

        for i := range cfg.Hosts {
            for v := range cfg.Hosts[i].Exclude {
               if (cfg.Hosts[i].Exclude[v] == collectorName) {
                    excluded = true
               }
          }
            if !excluded {
                list = append(list, IncludedHosts{cfg.Hosts[i].Ipaddress, cfg.Hosts[i].HostName,cfg.Hosts[i].Username, cfg.Hosts[i].Password,cfg.Hosts[i].RoutingEntryTime})
            }
        }
    return list
    }
