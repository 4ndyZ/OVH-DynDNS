package main

import (
	"encoding/csv"
	"errors"
	"net"
	"net/http"
	"strings"
)

func (a *App) needsRefresh() (bool, error) {
	needsRefresh := false
	// NS Lookup for domain
	ips, err := NSLookup(a.dynDNS.domain)
	if err != nil {
		return needsRefresh, err
	}
	// Check for no DNS entries at all
	if len(ips) == 0 {
		return true, nil
	}
	// Check if DNS entries still match
	for _, ipDNS := range ips {
		ipNow, err := a.getIP(GetIPType(ipDNS))
		if err != nil {
			break
		}
		if !ipNow.Equal(ipDNS) {
			needsRefresh = true
			break
		}
	}
	return needsRefresh, err
}

func (a *App) refresh() error {
	for _, ipType := range a.dynDNS.ipTypes {
		ip, err := a.getIP(ipType)
		if err != nil {
			return err
		}
		return a.ovh.Update(a.dynDNS.domain, ip)
	}
	return a.ovh.Refresh(a.dynDNS.domain)
}

func (a *App) getIP(ipType IPType) (net.IP, error) {
	// URL
	url := "https://api.andycraftz.eu/v1/ipv4"
	if ipType == IPv6 {
		url = "https://api.andycraftz.eu/v1/ipv6"
	}
	// Perform HTTP call
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Parse data
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	body, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	ipString := body[0][1]
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, errors.New(strings.Join([]string{ipString, " is not a valid IP address."}, ""))
	}
	return ip, nil
}
