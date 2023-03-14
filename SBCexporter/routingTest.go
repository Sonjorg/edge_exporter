//first request
package main
import (
	"encoding/xml"
	"fmt"
	//"log"
	//"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
)
type rt struct {
	//XMLname    xml.Name `xml:"_list"`
	Id         []string `xml:"_list>id,attr"`//`xml:"_pk,attr id="2" href="https://10.233.230.11/rest/routingtable//2"/>
	//<_pk id="4" href="https://10.233.230.11/rest/routingtable//4"/>
}
func test(){
	phpsessid, err :=  APISessionAuth("student", "panneKake23","https://10.233.230.11/rest/login")
	if err != nil {
	}
	data,err := getAPIData("https://10.233.230.11/rest/routingtable", phpsessid)
	if err != nil {
	}
	b := []byte(data) //Converting string of data to bytestream
	ssbc := &rt{}
	xml.Unmarshal(b, &ssbc) //Converting XML data to variables
	fmt.Println("Successful API call data: ",ssbc.Id,"\n")
}