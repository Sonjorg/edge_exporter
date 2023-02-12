package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type status struct {
    HTTPcode    int    `xml:"http_code"`
}

type SBCdata struct {
	XMLname    xml.Name  `xml:"root"`
	Status     status    `xml:"status"`
	//HTTPcode 				int      `xml:"status>http_code"`
	//Isdnsg					xml.Name `xml:"isdnsg"`
	//ActionsetTableNumber	int		 `xml:"isdnsg>ActionsetTableNumber"`
	//
}

func main() {
	data, _ := ioutil.ReadFile("APIoutput.xml")

	sbc := &SBCdata{}

	xml.Unmarshal([]byte(data), &sbc)

	fmt.Println(sbc.Status)
	//fmt.Println(note.From)
	//fmt.Println(note.Heading)
	//fmt.Println(note.Body)
}
