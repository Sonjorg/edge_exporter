package config

import (
	"fmt"
	"strings"

	//"os"
	"os"
	"strconv"
	//"flag"
)

//Describing docker-compose.yml file
type HostCompose struct {
    Authtimeout       int64  //"authtimeout"
    Hostname          string  //"hostname"
    Ip                string  //"ipaddress"
    Username          string  //"username"
    Password          string  //"password"
    Exclude         []string  //"exclude"
    RoutingEntryTime  float64 //"routing_database_hours"
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func GetConfig(hostCompose *HostCompose) *HostCompose{
        authtimeout, err := strconv.ParseInt(getEnv("Authtimeout", "2"), 10, 64)
        if err != nil {
			fmt.Print(err)
		}
        hostName := getEnv("hostname", "HostX")
        ipaddress := getEnv("ipaddress", "Ipaddress empty")
        username := getEnv("username", "Username empty")
        password := getEnv("password", "Password empty")
        excludeString := os.Getenv("exclude")
        exclude := strings.Split(excludeString, ",")
        routingEntry := getEnv("routing_database_hours", "24")
        routingEntryTime, err := strconv.ParseFloat(routingEntry,64)
        hostCompose = &HostCompose{
            Authtimeout: authtimeout, 
            Hostname: hostName, 
            Ip:       ipaddress, 
            Username: username, 
            Password: password, 
            Exclude: exclude, 
            RoutingEntryTime: routingEntryTime,
        }
    return hostCompose
}
func Excluded(collectorName string) bool {
    cfg := GetConfig(&HostCompose{})
    for i := range cfg.Exclude {
        if (cfg.Exclude[i] == collectorName) {
            return true
        }
    }
    return false
}
