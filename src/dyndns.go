package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

var srcurl = flag.String("srcurl", "", "url that supplies the current ip")
var token = flag.String("token", "", "token to authenticate against the dns update api")
var debug = flag.Bool("debug", false, "show debugging output")
var targetname = flag.String("target", "", "DNS A record to be updated")

func debugPrint(msg string) {
	if *debug == true {
		fmt.Printf("DEBUG: %s\n", msg)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func fetchIP(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	checkError(err)
	return string(b)
}

func lookupHost(hostname string) string {
	hostIP, err := net.LookupHost(hostname)
	checkError(err)
	return hostIP[0]
}

//func update()

func main() {
	flag.Parse()
	myIP := fetchIP(*srcurl)
	debugPrint(fmt.Sprintf("My public ip is %s", myIP))
	hostIP := lookupHost(*targetname)
	debugPrint(fmt.Sprintf("Host %s has ip %s", *targetname, hostIP))
	if myIP == hostIP {
		debugPrint(fmt.Sprintf("%s's ip is up to date, exiting", *targetname))
		os.Exit(0)
	}

}
