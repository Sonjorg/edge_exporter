//first request
package main
import (
	"encoding/xml"
	//"log"
	//"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
	"crypto/tls"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"log"
	"regexp"


)
type rt struct {
   // Value  float32 `xml:",chardata"`
	XMLName xml.Name `xml:"root"`
	Rt2     rt2      `xml:"routingtable_list"`
}
type rt2 struct {
	Rt3    rt3       `xml:"routingtable_pk"`
	//Value  float32 `xml:",chardata"`
	//Id         []int `xml:"id,attr"`//`xml:"_pk,attr id="2" href="https://10.233.230.11/rest/routingtable//2"/>
	//<_pk id="4" href="https://10.233.230.11/rest/routingtable//4"/>
}
type rt3 struct {
	Attr    []string `xml:"id,attr"`
	Value     string `xml:",chardata"`
}
 //Second request
 type call2xml1 struct {
	XMLName    xml.Name  `xml:"root"`
	Call2xml2  call2xml2 `xml:"routingentry_list"`
 }
 type call2xml2 struct {
	Call2xml3  call2xml3 `xml:"routingentry_pk"`
 }
 type call2xml3 struct {
	Attr    []string `xml:"id,attr"`
	Value     string `xml:",chardata"`
 }


//func routingCollector(ip string)([]prometheus.Metric, string){
	func main() {
	phpsessid, err := APISessionAuth("student", "PanneKake23","https://10.233.230.11/rest/login")
	if err != nil {
		//return nil, err.Error()
	}
	data,err := getAPIData("https://10.233.230.11/rest/routingtable", phpsessid)
	if err != nil {
		//return nil, err.Error()
	}
	b := []byte(data) //Converting string of data to bytestream
	ssbc := &rt{}
	xml.Unmarshal(b, &ssbc) //Converting XML data to variables
	//fmt.Println("Successful API call data: ",ssbc.Rt2.Rt3.Attr)
	routingTables := ssbc.Rt2.Rt3.Attr

	if (len(routingTables) <= 0) {
		//return nil, "Routingtables empty"
		fmt.Println("Routingtables empty")

	}
	for j := range routingTables {
    	url := "https://10.233.230.11/rest/routingtable/" + routingTables[j] + "/routingentry"
		data2, err := getAPIData(url, phpsessid)
		if err != nil {
		}
		b2 := []byte(data2) //Converting string of data to bytestream
		ssbc2 := &call2xml1{}
		xml.Unmarshal(b2, &ssbc2) //Converting XML data to variables
		routingEntries := ssbc2.Call2xml2.Call2xml3.Attr
		if (len(routingEntries)<=0) {
			continue
		}
		entries := regexp.MustCompile(`$\d+`)
		fmt.Println("Table:", routingEntries[j])

		for k := range routingEntries {
			match := entries.FindStringSubmatch(routingEntries[k])
			fmt.Println(match)
		}


	}

}



func APISessionAuth(username string, password string, loginURL string) (string,error) {
	//cfg := getConf(&Config{})
	//timeout := cfg.Authtimeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr,Timeout: 3 * time.Second}

	params := url.Values{}
	params.Add("Username", username)
	params.Add("Password", password)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", loginURL, body)
	if err != nil {
		log.Flags()
			fmt.Println("error in auth:", err)
			return "Error fetching data", err
		//	fmt.Println("error in systemExporter:", error)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Flags()
		fmt.Println("error in auth:", err)
		return "Error fetching data", err
		//fmt.Println("error in systemExporter:", err)
	}

	  m := make(map[string]string)
	  for _, c := range resp.Cookies() {
		 m[c.Name] = c.Value
	  }
	  fmt.Println(m["PHPSESSID"])
	  phpsessid := m["PHPSESSID"]

	defer resp.Body.Close()
	return phpsessid,err

}

func getAPIData(url string, phpsessid string) (string,error){

tr2 := &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
client2 := &http.Client{Transport: tr2}
cookie1 := &http.Cookie{
	Name:   "PHPSESSID",
	Value:  phpsessid,
	//Path:     "/",
	MaxAge:   3600,
	HttpOnly: false,
	Secure:   true,
}
req2, err := http.NewRequest("GET", url, nil)
if err != nil {
	log.Flags()
		fmt.Println("error in getapidata():", err)
		return "Error fetching data", err
	//	fmt.Println("error in systemExporter:", error)
}
req2.AddCookie(cookie1)
	resp2, err := client2.Do(req2)
	if err != nil {
		log.Flags()
			fmt.Println("error in getapidata():", err)
			return "Error fetching data", err
	}

	b, err := ioutil.ReadAll(resp2.Body)
	defer resp2.Body.Close()

	return string(b), err
}
