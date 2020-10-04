package dyndns

import "testing"

func TestFetchIP(t *testing.T) {
	testURL := "https://www.drexler-online.net/testmyip"
	myIP := fetchIP(testURL)
	t.Logf("Got ip %s from %s", myIP, testURL)
	testIP := "10.10.10.10"
	if myIP != "10.10.10.10" {
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
