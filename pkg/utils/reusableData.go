package utils

import "fmt"

//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
//_ "github.com/mattn/go-sqlite3"
//"fmt"

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

	m :=  make(map[string][]string)
	m2 :=  make(map[string]RoutingTablesUtils)

	s := RoutingData{}
	t := RoutingTablesUtils{}
	t.Time = time

	m[routingTable] = routingEntries

		t.RoutingEntries = m

	t.RoutingTables = append(t.RoutingTables, routingTable)
	m2[ipaddress] = t
	s.Routing = m2

}

func GetRoutingData(ipaddress string, r RoutingData) (map[string][]string,[]string,string) {

	fmt.Println(r.Routing[ipaddress].RoutingEntries, r.Routing[ipaddress].RoutingTables,r.Routing[ipaddress].Time)

	_, ok := r.Routing[ipaddress]
	if !ok {
		return nil,nil,	"no data"
	}
	return r.Routing[ipaddress].RoutingEntries, r.Routing[ipaddress].RoutingTables,r.Routing[ipaddress].Time

}