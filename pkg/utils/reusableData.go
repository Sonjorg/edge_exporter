package utils

import (
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//_ "github.com/mattn/go-sqlite3"
	//"fmt"
)

type RoutingTablesUtils struct {
	Id                int
	Time              string
	RoutingTables   []string
	RoutingEntries  map[string][]string

	 //map consisting of routingtables and their routingentries
}
type RoutingData struct {
	Routing  map[string]RoutingTablesUtils
}
func StoreRoutingEntries(ipaddress string, time string, routingTable string, routingEntries []string) {

	s := RoutingData{}
	t := RoutingTablesUtils{}
	t.Time = time
	if t.RoutingEntries[routingTable] == nil {

		t.RoutingEntries[routingTable] = routingEntries
	}
	t.RoutingTables = append(t.RoutingTables, routingTable)
	s.Routing[ipaddress] = t

}

func GetRoutingData(ipaddress string, r RoutingData) (map[string][]string,[]string,string) {

	_, ok := r.Routing[ipaddress]
	if !ok {
		return nil,nil,	"no data"
	}
	return r.Routing[ipaddress].RoutingEntries, r.Routing[ipaddress].RoutingTables,r.Routing[ipaddress].Time

}