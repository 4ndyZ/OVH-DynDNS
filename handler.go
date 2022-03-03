package main

import (
	"encoding/csv"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

func (a *App) needsRefresh() (bool, error) {
	// Always perform refresh when check methode is API
	if a.dynDNS.checkMethod == API {
		return true, nil
	}
	// Get set IPs
	var ips []net.IP
	if ips == nil || len(ips) == 0 {
		return true, nil
	}
	//
	needsRefresh := false
	// Check if entries still match
	for _, ipSet := range ips {
		ipNow, err := a.getIP(GetIPType(ipSet))
		if err != nil {
			break
		}
		if !ipNow.Equal(ipSet) {
			needsRefresh = true
			break
		}
	}
	return needsRefresh, nil
}

func (a *App) refresh() error {
	for _, ipType := range a.dynDNS.ipTypes {
		ip, err := a.getIP(ipType)
		if err != nil {
			return err
		}
		return a.ovh.Update(a.dynDNS.domain, ip, false)
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
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			Log.Logger.Warn().Str("error", err.Error()).Msg("Error while closing HTTP call response body.")
		}
	}(resp.Body)
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
