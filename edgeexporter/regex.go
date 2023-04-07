// Read this repo to understand better: https://github.com/antonputra/tutorials/blob/main/lessons/141/prometheus-nginx-exporter/cmd/exporter/main.go
// So far it is scraping one metric from a txt file
// Edit path of file
package main

import (
	//	"bytes"
	//	"flag"
	//	"fmt"
	//	"io"
	"log"
	//	"net/http"
	"regexp"
	//	"strconv"
	//	"time"
	"fmt"
	//exporter "github.com/antonputra/tutorials/lessons/141/prometheus-nginx-exporter"
	"github.com/hpcloud/tail"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

// [<ChannelClearingDelay>][\d][</ChannelClearingDelay>]
// var exp = regexp.MustCompile(`^(?P<remote>[^ ]*) (?P<host>[^ ]*)$`)
var exp = regexp.MustCompile(`<ChannelClearingDelay>\d+</ChannelClearingDelay>`)
var value = regexp.MustCompile(`\d+`)

//regexp.MustCompile(`^(?P<remote>[^ ]*) (?P<host>[^ ]*) (?P<user>[^ ]*) \[(?P<time>[^\]]*)\] \"(?P<method>\w+)(?:\s+(?P<path>[^\"]*?)(?:\s+\S*)?)?\" (?P<status_code>[^ ]*) (?P<size>[^ ]*)(?:\s"(?P<referer>[^\"]*)") "(?P<agent>[^\"]*)" (?P<urt>[^ ]*)$`)

func regex() {
	t, err := tail.TailFile("C:/Users/sonde/Desktop/HDOmonitoring/SBCexporter/tailing-test.txt", tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Fatalf("tail.TailFile failed: %s", err)
	}

	fmt.Println(t, err)

	for line := range t.Lines {
		matchType := exp.FindStringSubmatch(line.Text)
		matchValue := value.FindStringSubmatch(line.Text)

		fmt.Println(matchType, ": ", matchValue)

		//result := make(map[string]string)

	}

}
