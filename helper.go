package main

import (
	"golang.org/x/net/publicsuffix"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// IsIP check if input string is a IP and returns a boolean
func IsIP(ipString string) bool {
	ip := net.ParseIP(ipString)
	return ip != nil
}

// GetIPType return the IPType of if a ip (net.IP) object (IPv4 oder IPv6)
func GetIPType(ip net.IP) IPType {
	if ip.To4() != nil {
		return IPv4
	}
	return IPv6
}

// IsDomain checks if a string is a domain and returns a boolean
func IsDomain(domain string) bool {
	if len(domain) == 0 {
		return false
	}
	if IsIP(domain) {
		return false
	}
	if len(strings.Replace(domain, ".", "", -1)) > 253 {
		return false
	}
	regex := regexp.MustCompile(`\b((xn--)?[a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}\b`)
	return regex.MatchString(domain)
}

/*
GetZoneFromDomain exctracts the root domain zone as string from a domain string

For example:
 - Parameter: "subdomain.example.com"
 - Return: "example.com"
*/
func GetZoneFromDomain(domain string) string {
	zone, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return ""
	}
	return zone
}

/*
GetSubDomainFromDomain extracts the subdomain as string from a domain string

For example:
 - Parameter: "subdomain.example.com"
 - Return: "subdomain"
*/
func GetSubDomainFromDomain(domain string) string {
	zone := GetZoneFromDomain(domain)
	if len(domain) == len(zone) {
		return ""
	}
	return domain[:len(domain)-len(zone)-1] // -1 to cut away the dot
}

/*
NSLookup performs a DNS lookup for a domain string and return all IPs from the A or AAAA records as an array
*/
func NSLookup(domain string) ([]net.IP, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		// Check if error is not found and return empty array
		if err.(*net.DNSError).IsNotFound {
			return []net.IP{}, nil
		}
		return []net.IP{}, err
	}
	return ips, nil
}

// IntToString converts a integer value to an string
func IntToString(i int) string {
	return strconv.Itoa(i)
}
