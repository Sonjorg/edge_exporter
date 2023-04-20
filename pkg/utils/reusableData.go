package utils

import "fmt"

//import "fmt"

//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
//_ "github.com/mattn/go-sqlite3"
//"fmt"

type RoutingTablesUtils struct {
	Id                int
	Time              string
	RoutingTables   []string
	Ipaddress string
	RoutingEntries  map[string][]string

	 //map consisting of routingtables and their routingentries
}
type RoutingData struct {
	Ipaddress string
	Routing  RoutingTablesUtils
}
func StoreRoutingEntries(ipaddress string, time string, routingTable string, routingEntries []string) {

	m :=  make(map[string][]string)
	m2 :=  make(map[string]RoutingTablesUtils)

	s := []RoutingTablesUtils{}
	t := RoutingTablesUtils{}
	t.Time = time

	m[routingTable] = routingEntries

	a := RoutingTablesUtils{Ipaddress: ipaddress, Time: time, RoutingEntries: m}
	a.RoutingTables = append(a.RoutingTables, routingTable)
	t.RoutingEntries = m

	t.RoutingTables = append(t.RoutingTables, routingTable)
	m2[ipaddress] = t
	//Ipaddress := ipaddress

		s = append(s, a)
		//	s.Routing = append(s.Routing,t)

}

func GetRoutingData(ipaddress string, r RoutingData) (map[string][]string,[]string,string) {

	//fmt.Println(r.Routing[ipaddress].RoutingEntries, r.Routing[ipaddress].RoutingTables,r.Routing[ipaddress].Time)

	/*_, ok := r.Routing[ipaddress]
	if !ok {
		return nil,nil,	"no data"
	}*/
	s := []RoutingTablesUtils{}
	for i:= range s {
		if (s[i].Ipaddress == ipaddress) {
			fmt.Println(s[i].RoutingEntries, s[i].RoutingTables,s[i].Time)
			return s[i].RoutingEntries, s[i].RoutingTables,s[i].Time
		}

	}

	return nil,nil,"no data"
}