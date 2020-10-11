package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numerals = "0123456789"
)

func StartTestServer(retval string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, retval)
	}))

	return ts
}

func GenerateRandomString(length int) string {
	pool := Alphabet + Numerals
	rand.Seed(time.Now().UnixNano())
	var ret []string
	//fmt.Print(pool)
	i := 0
	for ; i < length; i++ {
		rnum := rand.Intn(len(pool))
		rchar := pool[rnum]
		//fmt.Printf("%d: %s\n", rnum, string(rchar))
		ret = append(ret, string(rchar))
	}
	return strings.Join(ret, "")
}

func TestFetchIP(t *testing.T) {
	randString := GenerateRandomString(12)
	t.Logf("Expected string: %s", randString)
	ts := StartTestServer(randString)
	defer ts.Close()
	testURL := ts.URL
	myIP := fetchIP(testURL)
	t.Logf("Got ip %s from %s", myIP, testURL)
	testIP := randString
	if myIP != randString {
		t.Errorf("Testing ip wrong: got %s, expected: %s", myIP, testIP)
	}
}


func TestLookupHost(t *testing.T) {
	testHostName := "www.drexler-online.net"
	expectedIP := "95.216.59.146"
	testIP := lookupHost(testHostName)
	t.Logf("looking up A record for %s, got %s", testHostName, testIP)
	if testIP != expectedIP {
		t.Errorf("Got %s, expected %s for host %s", testIP, expectedIP, testHostName)
	}
}

func TestSplitDomainName(t *testing.T) {
	testHostName := "www.drexler-online.net"
	expectedDomain := "drexler-online.net"
	expectedHostName := "www"
	hostName := splitDomainName(testHostName)[0]
	domain := splitDomainName(testHostName)[1]
	if domain != expectedDomain {
		t.Errorf("Got %s, expected %s", domain, expectedDomain)
	}
	if hostName != expectedHostName {
		t.Errorf("Got %s, expected %s", hostName, expectedHostName)
	}
}
