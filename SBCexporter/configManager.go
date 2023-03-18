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
            //exclude        string `yaml:"exclude"`
                    // Server is the general server timeout to use
                    // for graceful shutdowns
                Exclude struct {
                        // Server is the general server timeout to use
                        // for graceful shutdowns
                      SystemExporter bool `yaml:"systemstats"`
                     routingEntry   bool `yaml:"routingentry"`
                }`yaml:"exclude"`

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

     func getIncludedHosts(exporterName string) []includedHosts{
        cfg := getConf(&Config{})

        //var list []includedHosts
        list := make([]includedHosts,0,8)
        switch exporterName {
            case "systemStats":
               for i := range cfg.Hosts {
                    if (cfg.Hosts[i].Exclude.SystemExporter == false) {
                        list = append(list, includedHosts{cfg.Hosts[i].Ipaddress, cfg.Hosts[i].HostName,cfg.Hosts[i].Username, cfg.Hosts[i].Password})
                    }
                }
            case "callStats":
                for i:= range cfg.Hosts {
                    if (cfg.Hosts[i].Exclude.routingEntry == false) {
                        list = append(list, includedHosts{cfg.Hosts[i].Ipaddress, cfg.Hosts[i].HostName,cfg.Hosts[i].Username, cfg.Hosts[i].Password})
                    }
                }
                //INFO: have a switch case on all exporters made, must remember exact exporternames inside each exporter
            }
    return list
    }

func getAuth(ipadr string) (username string, password string) {
    var u, p string
    cfg := getConf(&Config{})

    for i:= range cfg.Hosts {
        if (cfg.Hosts[i].Ipaddress == ipadr) {
            u, p = cfg.Hosts[i].Username, cfg.Hosts[i].Password
        }
    }
   // return "test", "test"
    return u,p
}
//func IndexFunc[E any](s []E, f func(E) bool) int
func getHostName(ipaddress string) string{
    cfg := getConf(&Config{})
    var host string
    for i := range cfg.Hosts {
        if (cfg.Hosts[i].Ipaddress == ipaddress) {
            host = cfg.Hosts[i].HostName
        }
    }
    return host
}
/*
func main() {
   // ip := getIpAdrExp("systemStats")
    //fmt.Println(ip)
   // conf := getConf(&Config{})
    //conf.Hosts.Exclude
    //v:= getHostName("46.333.534.22")
    g:= getIncludedHosts("systemStats")
    for i:= range g {
    fmt.Println(g[i].hostname,g[i].username)
    }
}*/