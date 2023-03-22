package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"strconv"

	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// The functions APISessionAuth(...) and getAPIData(...) utilizes curl-to-go translator but is modified for cookie management and storing these with sqlite.
// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

// curl -k --data "Username=student&Password=PanneKake23" -i -v https://10.233.230.11/rest/login

// TODO: This is insecure; use only in dev environments.
func APISessionAuth(username string, password string, ipaddress string) (string,error) {
	//var read []byte
	var phpsessid string
	var err error
	phpsessid,err = getSqliteData(ipaddress)
	//fmt.Println(phpsessid)
	if (phpsessid != "") {
		fmt.Println("henta fra sql",err)

		return phpsessid, nil
	}


	cfg := getConf(&Config{})
	timeout := cfg.Authtimeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr,Timeout: time.Duration(timeout) * time.Second}

	params := url.Values{}
	params.Add("Username", username)
	params.Add("Password", password)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://"+ipaddress+"/rest/login", body)
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
	phpsessid = m["PHPSESSID"]


	defer resp.Body.Close()


	//fmt.Println("henta fra ruter")
	//os.Remove("sqlite-database.db")
	//Checking if db already exist
	_, err = os.Stat("sqlite-database.db")
	if err != nil {
		fmt.Println("Creating sqlite-database.db...")
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			fmt.Println(err.Error())
		}
		file.Close()
	}

	var sqliteDatabase *sql.DB

	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		fmt.Println(err)
	}

	err = createTable(sqliteDatabase)
	if err != nil {
		fmt.Println(err)
	}
	now := time.Now().Format(time.RFC3339)

	if rowExists(sqliteDatabase,ipaddress) {
		Update(sqliteDatabase, phpsessid, now, ipaddress)
		//insertAuth(sqliteDatabase, ipaddress, phpsessid, time.Now().String())
	} else {
		insertAuth(sqliteDatabase, ipaddress, phpsessid, now)
	}

	return phpsessid,nil
	}


//This function retrieves the session cookie from the sqlite database
func getSqliteData(ipaddress string) (cookie string, err error){

	var sqliteDatabase *sql.DB

	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		fmt.Println(err)
		return "", err
	} // Open the created SQLite File
	 // Defer Closing the database


	Hosts,err := displayAuth(sqliteDatabase,ipaddress)
	if err != nil {
		return "", err
	}
	defer sqliteDatabase.Close()
	//defer file.Close()
	var c string
	mins := time.Minute * time.Duration(8)
	for i:= range Hosts {
		fmt.Println(Hosts[i].Ipaddress)
		if (Hosts[i].Ipaddress == ipaddress) {

			now := time.Now().Format(time.RFC3339)
			//previous := time.
			parsed, _ := time.Parse(time.RFC3339, now)
			parsed2,err := time.Parse(time.RFC3339, Hosts[i].Time)
			if err != nil {
				fmt.Println(err)
				return "",nil
			}
			//fmt.Println(parsed2.Add(mins).After(parsed))

			if (parsed2.Add(mins).After(parsed) == true) {
				//return Hosts[i].Phpsessid,nil
				//phpsessid = Hosts[i].Phpsessid
				//fmt.Println("true")
				c = Hosts[i].Phpsessid
				//fmt.Println(c)
				break

			} else {
				 return "", nil}
		} else {
			return "", nil}
	}
	return c, nil
}
/*
func main() {
	/*var sqliteDatabase *sql.DB

	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		fmt.Println("cant open db")
	} // Open the created SQLite File
	 // Defer Closing the database
	// insertAuth(sqliteDatabase, "10.233.234.11", "test 19.00", time.Now().String())
	 s,err  := APISessionAuth("student","PanneKake23", "10.233.234.11")
	 if err != nil {
		fmt.Println(err)
	 }
	 fmt.Println(s)


}
*/
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