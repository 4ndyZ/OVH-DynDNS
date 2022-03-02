package main

import (
	"errors"
	"github.com/ovh/go-ovh/ovh"
	"net"
	"strings"
)

// OVH struct to hold API client
type OVH struct {
	client *ovh.Client
}

// Initialize OVH struct with configuration
func (o *OVH) Initialize(configuration Configuration) error {
	client, err := ovh.NewClient(
		configuration.OVH.Region,
		configuration.OVH.ApplicationKey,
		configuration.OVH.ApplicationSecretKey,
		configuration.OVH.ConsumerKey,
	)
	if err != nil {
		return err
	}
	o.client = client
	return nil
}

// Update DNS Record by using domain and new IP using the OVH API
func (o *OVH) Update(domain string, ip net.IP) error {
	recordID, err := o.getRecordID(domain, ip)
	if err != nil {
		return err
	}
	if recordID == 0 {
		return o.createRecord(domain, ip)
	}
	return o.updateRecord(domain, ip, recordID)
}

// Refresh saves the changed DNS zone
func (o *OVH) Refresh(domain string) error {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/refresh/"}, "")
	// TODO: DEBUG endpoint
	return o.client.Post(endpoint, nil, nil)
}

// getRecordID get the A or AAAA record id of a domain from the OVH API
func (o *OVH) getRecordID(domain string, ip net.IP) (int, error) {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/record?fieldType=", GetIPType(ip).String(), "&subDomain=", GetSubDomainFromDomain(domain)}, "")
	// TODO: DEBUG endpoint
	var domains []int
	err := o.client.Get(endpoint, &domains)
	if err != nil {
		return 0, err
	}
	// TODO: DEBUG domains
	if len(domains) == 0 {
		return 0, nil
	}
	if len(domains) > 1 {
		return 0, errors.New(strings.Join([]string{"Found", IntToString(len(domains)), "entries for domain ", domain, "."}, ""))
	}
	return domains[0], nil
}

// createRecord create a new A or AAAA record for a defined domain using the OVH API
func (o *OVH) createRecord(domain string, ip net.IP) error {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/record"}, "")
	params := createRecordParams{
		IPType:    GetIPType(ip).String(),
		SubDomain: GetSubDomainFromDomain(domain),
		Target:    ip.String(),
		TTL:       60,
	}
	// TODO: DEBUG endpoint + params
	return o.client.Post(endpoint, &params, nil)
}

// updateRecord update an A or AAAA record for a defined domain using the OVH API
func (o *OVH) updateRecord(domain string, ip net.IP, recordID int) error {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/record/", IntToString(recordID)}, "")
	params := updateRecordParams{
		SubDomain: GetSubDomainFromDomain(domain),
		Target:    ip.String(),
		TTL:       60,
	}
	// TODO: DEBUG endpoint + params
	return o.client.Put(endpoint, &params, nil)
}

// Structs

type createRecordParams struct {
	IPType    string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}

type updateRecordParams struct {
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}
