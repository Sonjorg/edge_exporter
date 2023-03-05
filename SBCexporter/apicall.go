package main

import (
	"crypto/tls"
	"fmt"
	"strings"
    //"bufio"
	//"io/ioutil"
	//"github.com/tiket-oss/phpsessgo"
	"io/ioutil"
	"net/http"
//"net/http/cookiejar"
	"net/url"
	//"regexp"
	//"strconv"
	"log"
)

//The functions APISessionAuth(...) and getAPIData(...) utilizes curl-to-go translator but is modified: https://mholt.github.io/curl-to-go/
// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

// curl -k --data "Username=student&Password=PanneKake23" -i -v https://10.233.230.11/rest/login

// TODO: This is insecure; use only in dev environments.
func APISessionAuth(username string, password string, loginURL string) (string,error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	params := url.Values{}
	params.Add("Username", username)
	params.Add("Password", password)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", loginURL, body)
	if err != nil {
		log.Flags()
			fmt.Println("error in auth:", err)
		//	fmt.Println("error in systemExporter:", error)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Flags()
		fmt.Println("error in auth:", err)
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
//The following two comments are produced by the curl-to-go utility: https://mholt.github.io/curl-to-go/
//The code produced by curl-to-go has added functionality such as sending cookies

// curl --cookie  \ -i -k https://10.233.230.11/rest/isdnsg/10001

// TODO: This is insecure; use only in dev environments.

func getAPIData(url string, phpsessid string) (string){

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
		fmt.Println("error in http request:", err)
}
req2.AddCookie(cookie1)
	resp2, err := client2.Do(req2)
	if err != nil {
		fmt.Println("error in http request:", err)
}

	b, err := ioutil.ReadAll(resp2.Body)

	defer resp2.Body.Close()
	return string(b)
}

