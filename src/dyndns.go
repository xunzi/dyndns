package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var srcurl = flag.String("srcurl", "", "url that supplies the current ip")
var token = flag.String("token", "", "token to authenticate against the dns update api")
var debug = flag.Bool("debug", false, "show debugging output")
var targetname = flag.String("target", "", "DNS A record to be updated")

var HetznerDnsAPI = "https://dns.hetzner.com/api/v1"
var ApiToken = os.Getenv("DNSTOKEN")

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

func hetzerFetchZoneID(domainname string) string {
	client := &http.Client{}
	type resultMap struct {
		Zones []struct {
			ID            string   `json:"id"`
			Name          string   `json:"name"`
			TTL           int      `json:"ttl"`
			Registrar     string   `json:"registrar"`
			LegacyDNSHost string   `json:"legacy_dns_host"`
			LegacyNs      []string `json:"legacy_ns"`
			Ns            []string `json:"ns"`
			Created       string   `json:"created"`
			Verified      string   `json:"verified"`
			Modified      string   `json:"modified"`
			Project       string   `json:"project"`
			Owner         string   `json:"owner"`
			Permission    string   `json:"permission"`
			ZoneType      struct {
				ID          string      `json:"id"`
				Name        string      `json:"name"`
				Description string      `json:"description"`
				Prices      interface{} `json:"prices"`
			} `json:"zone_type"`
			Status          string `json:"status"`
			Paused          bool   `json:"paused"`
			IsSecondaryDNS  bool   `json:"is_secondary_dns"`
			TxtVerification struct {
				Name  string `json:"name"`
				Token string `json:"token"`
			} `json:"txt_verification"`
			RecordsCount int `json:"records_count"`
		} `json:"zones"`
		Meta struct {
			Pagination struct {
				Page         int `json:"page"`
				PerPage      int `json:"per_page"`
				PreviousPage int `json:"previous_page"`
				NextPage     int `json:"next_page"`
				LastPage     int `json:"last_page"`
				TotalEntries int `json:"total_entries"`
			} `json:"pagination"`
		} `json:"meta"`
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zones", HetznerDnsAPI), nil)
	checkError(err)
	req.Header.Set("Auth-API-Token", ApiToken)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("search_name", domainname)
	req.URL.RawQuery = q.Encode()
	debugPrint(req.URL.String())
	resp, err := client.Do(req)
	checkError(err)
	debugPrint(resp.Status)
	if resp.StatusCode != 200 {
		log.Fatalf("Request for %s returned %s", q.Encode(), resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	checkError(err)
	var zoneResult resultMap
	json.Unmarshal(b, &zoneResult)
	//debugPrint(fmt.Sprintf("%v", resultMap))
	return zoneResult.Zones[0].ID
}

func checkHTTPStatus(status int, statuscode string)  {
	if status != 200 {
		log.Fatalf("Http status code is %s", statuscode)
	}
}

func splitDomainName(hostname string) []string {
	splitName := strings.SplitN(hostname, ".", 2)
	return splitName
}

func hetzerFetchRecordID(hostname string, zoneid string) string {
	client := &http.Client{}
	type resultMap struct {
		Records []struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Name     string `json:"name"`
			Value    string `json:"value"`
			TTL      int    `json:"ttl,omitempty"`
			ZoneID   string `json:"zone_id"`
			Created  string `json:"created"`
			Modified string `json:"modified"`
		} `json:"records"`
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/records", HetznerDnsAPI), nil)
	checkError(err)
	req.Header.Set("Auth-API-Token", ApiToken)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("zone_id", zoneid)
	req.URL.RawQuery = q.Encode()
	debugPrint(req.URL.String())
	resp, err := client.Do(req)
	checkError(err)
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	checkError(err)
	//debugPrint(string(b))
	var recordResult resultMap
	json.Unmarshal(b, &recordResult)
	for _, record := range recordResult.Records {
		if record.Name == hostname && record.Type == "A" {
			return record.ID
		}
	}
	return ""
}

func hetznerUpdateDNSRecord(recordid string, name string, ip string, zoneid string) {
	client := &http.Client{}
	type updateRecord struct {
		Zone_id string
		Name    string
		Type    string
		Value   string
	}
	newRecord := updateRecord{
		Zone_id: zoneid,
		Name:    name,
		Type:    "A",
		Value:   ip,
	}
	var jsonData []byte
	jsonData, err := json.Marshal(newRecord)
	checkError(err)
	body := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/records/%s", HetznerDnsAPI, recordid), body)
	checkError(err)
	req.Header.Set("Auth-API-Token", ApiToken)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	checkError(err)
	if resp.StatusCode != 200 {
		log.Fatalf("Request for %s returned status %s", req.URL, resp.Status)
	}
	debugPrint(fmt.Sprintf("DNS entry update request returned status %s", resp.Status))

}

func main() {
	flag.Parse()
	if *targetname == "" {
		log.Fatal("Please supply a targetname as argument")
		}
	myIP := fetchIP(*srcurl)
	debugPrint(fmt.Sprintf("My public ip is %s", myIP))
	hostIP := lookupHost(*targetname)
	debugPrint(fmt.Sprintf("Host %s has ip %s", *targetname, hostIP))
	if myIP == hostIP {
		debugPrint(fmt.Sprintf("%s's ip is up to date, exiting", *targetname))
		os.Exit(0)
	}

	if ApiToken == "" {
		log.Fatal("Please supply a valid API token as environment var DNSTOKEN, e.g. export DNSTOKEN=123456xyz")
	}
	hostPart := splitDomainName(*targetname)[0]
	domain := splitDomainName(*targetname)[1]
	debugPrint(hostPart)
	debugPrint(domain)
	zoneID := hetzerFetchZoneID(domain)
	debugPrint(fmt.Sprintf("Zoneid for zone %s is %s", domain, zoneID))
	recordID := hetzerFetchRecordID(hostPart, zoneID)
	debugPrint(fmt.Sprintf("DNS entry %s in zone %s has record id %s", hostPart, domain, recordID))
	hetznerUpdateDNSRecord(recordID, hostPart, myIP, zoneID)
}
