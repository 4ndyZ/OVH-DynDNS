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
func (o *OVH) Update(domain string, ip net.IP, force bool) error {
	recordID, err := o.getRecordID(domain, ip)
	if err != nil {
		return err
	}
	// Create record when recordID does not exist
	if recordID == 0 {
		return o.createRecord(domain, ip)
	}
	if !force {
		// Get record for checking
		record, err := o.getRecord(domain, recordID)
		if err != nil {
			return err
		}
		// Check if DNS record is already correct
		if ip.String() == record.Target {
			return nil
		}
	}
	return o.updateRecord(domain, ip, recordID)
}

// Refresh saves the changed DNS zone
func (o *OVH) Refresh(domain string) error {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/refresh/"}, "")
	Log.Logger.Debug().Str("method", "POST").Str("endpoint", endpoint).Msg("Update DNS Zone.")
	return o.client.Post(endpoint, nil, nil)
}

// getRecordID get the A or AAAA record id of a domain from the OVH API
func (o *OVH) getRecordID(domain string, ip net.IP) (int, error) {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/record?fieldType=", GetIPType(ip).String(), "&subDomain=", GetSubDomainFromDomain(domain)}, "")
	var domains []int
	Log.Logger.Debug().Str("method", "GET").Str("endpoint", endpoint).Str("subdomain", GetSubDomainFromDomain(domain)).Str("record-type", GetIPType(ip).String()).Msg("Get DNS record.")
	err := o.client.Get(endpoint, &domains)
	if err != nil {
		return 0, err
	}
	if len(domains) == 0 {
		return 0, nil
	}
	if len(domains) > 1 {
		return 0, errors.New(strings.Join([]string{"Found", IntToString(len(domains)), "entries for domain ", domain, "."}, ""))
	}
	return domains[0], nil
}

func (o *OVH) getRecord(domain string, recordID int) (getRecord, error) {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/record/", IntToString(recordID)}, "")
	var record getRecord
	Log.Logger.Debug().Str("method", "GET").Str("endpoint", endpoint).Str("subdomain", GetSubDomainFromDomain(domain)).Int("record-id", recordID).Msg("Get DNS record ID.")
	err := o.client.Get(endpoint, &record)
	if err != nil {
		return record, err
	}
	return record, nil
}

// createRecord create a new A or AAAA record for a defined domain using the OVH API
func (o *OVH) createRecord(domain string, ip net.IP) error {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/record"}, "")
	record := createRecord{
		IPType:    GetIPType(ip).String(),
		SubDomain: GetSubDomainFromDomain(domain),
		Target:    ip.String(),
		TTL:       60,
	}
	Log.Logger.Debug().Str("method", "POST").Str("endpoint", endpoint).Str("subdomain", GetSubDomainFromDomain(domain)).Str("target", ip.String()).Str("record-type", GetIPType(ip).String()).Msg("Create DNS record.")
	return o.client.Post(endpoint, &record, nil)
}

// updateRecord update an A or AAAA record for a defined domain using the OVH API
func (o *OVH) updateRecord(domain string, ip net.IP, recordID int) error {
	endpoint := strings.Join([]string{"/domain/zone/", GetZoneFromDomain(domain), "/record/", IntToString(recordID)}, "")
	record := updateRecord{
		SubDomain: GetSubDomainFromDomain(domain),
		Target:    ip.String(),
		TTL:       60,
	}
	Log.Logger.Debug().Str("method", "PUT").Str("endpoint", endpoint).Str("subdomain", GetSubDomainFromDomain(domain)).Str("target", ip.String()).Str("record-type", GetIPType(ip).String()).Msg("Update DNS record.")
	return o.client.Put(endpoint, &record, nil)
}

// getRecord for holding refs for OVH DNS record
type getRecord struct {
	ID        int    `json:"id"`
	IPType    string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	Zone      string `json:"zone"`
	TTL       int    `json:"ttl"`
}

// getRecord for holding refs for OVH DNS record for creation
type createRecord struct {
	IPType    string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}

// getRecord for holding refs for OVH DNS record for update
type updateRecord struct {
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}
