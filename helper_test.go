package main

import (
	"net"
	"reflect"
	"sort"
	"testing"
)

func Test_IsIP(t *testing.T) {
	// Sample data
	ipv4 := "1.1.1.1"
	ipv6 := "2606:4700:4700::1111"
	domain := "example.com"
	noIP := "900.33.1.1"

	result := IsIP(ipv4)
	if result != true {
		t.Errorf("Got wrong result for valid IP %s", ipv4)
	}

	result = IsIP(ipv6)
	if result != true {
		t.Errorf("Got wrong result for valid IP %s", ipv6)
	}

	result = IsIP(domain)
	if result != false {
		t.Errorf("Got wrong result for non IP %s", domain)
	}

	result = IsIP(noIP)
	if result != false {
		t.Errorf("Got wrong result for invalid IP %s", noIP)
	}
}

func Test_GetIPType(t *testing.T) {
	ipv4 := net.ParseIP("1.1.1.1")
	ipv6 := net.ParseIP("2606:4700:4700::1111")

	result := GetIPType(ipv4)
	if result != IPv4 {
		t.Errorf("Got wrong result for IPv4 %s", ipv4.String())
	}

	result = GetIPType(ipv6)
	if result != IPv6 {
		t.Errorf("Got wrong result for IPv6 %s", ipv6.String())
	}
}

func Test_IsDomain(t *testing.T) {
	// Sample data
	ipv4 := "1.1.1.1"
	ipv6 := "2606:4700:4700::1111"
	domain := "example.com"
	noDomain := "com"

	result := IsDomain(ipv4)
	if result != false {
		t.Errorf("Got wrong result for IP %s", ipv4)
	}

	result = IsDomain(ipv6)
	if result != false {
		t.Errorf("Got wrong result for IP %s", ipv6)
	}

	result = IsDomain(domain)
	if result != true {
		t.Errorf("Got wrong result for valid domain %s", domain)
	}

	result = IsDomain(noDomain)
	if result != false {
		t.Errorf("Got wrong result for invalid domain %s", noDomain)
	}
}

func Test_GetZoneFromDomain(t *testing.T) {
	// Sample data
	domainWithSubDomain := "sub.example.co.uk"
	domainWithSubDomain2 := "subsub.sub.example.bayern"
	domainWithoutSubDomain := "example.com"

	result := GetZoneFromDomain(domainWithSubDomain)
	if result != "example.co.uk" {
		t.Errorf("Got wrong zone %s for domain %s", result, domainWithSubDomain)
	}

	result = GetZoneFromDomain(domainWithSubDomain2)
	if result != "example.bayern" {
		t.Errorf("Got wrong zone %s for domain %s", result, domainWithSubDomain2)
	}

	result = GetZoneFromDomain(domainWithoutSubDomain)
	if result != "example.com" {
		t.Errorf("Got wrong zone %s for domain %s", result, domainWithoutSubDomain)
	}
}

func Test_GetSubDomainFromDomain(t *testing.T) {
	// Sample data
	domainWithSubDomain := "sub.example.co.uk"
	domainWithSubDomain2 := "subsub.sub.example.bayern"
	domainWithoutSubDomain := "example.com"

	result := GetSubDomainFromDomain(domainWithSubDomain)
	if result != "sub" {
		t.Errorf("Got wrong subdomain %s for domain %s", result, domainWithSubDomain)
	}

	result = GetSubDomainFromDomain(domainWithSubDomain2)
	if result != "subsub.sub" {
		t.Errorf("Got wrong subdomain %s for domain %s", result, domainWithSubDomain2)
	}

	result = GetSubDomainFromDomain(domainWithoutSubDomain)
	if result != "" {
		t.Errorf("Got wrong subdomain %s for domain %s", result, domainWithoutSubDomain)
	}
}

func Test_NSLookup(t *testing.T) {
	// Sample data
	domain := "one.one.one.one"
	ipsExpected := []string{
		"1.1.1.1",
		"1.0.0.1",
		"2606:4700:4700::1001",
		"2606:4700:4700::1111",
	}
	sort.Strings(ipsExpected)
	// Check
	ips, err := NSLookup(domain)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	ipsIs := []string{}
	for _, ipDNS := range ips {
		ipsIs = append(ipsIs, ipDNS.String())
	}
	sort.Strings(ipsIs)
	// Compare is to expected array
	if !reflect.DeepEqual(ipsIs, ipsExpected) {
		t.Errorf("Got wrong IPs %v for domain %s expected %v", ipsIs, domain, ipsExpected)
	}
}

func Test_IntToString(t *testing.T) {
	result := IntToString(1337)
	if result != "1337" {
		t.Errorf("Got wrong value %s for int %s", result, "1337")
	}
}
