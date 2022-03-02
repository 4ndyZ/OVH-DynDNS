package main

import (
	"strings"
)

// App struct to hold refs
type App struct {
	ovh    OVH
	dynDNS struct {
		domain  string
		ipTypes []IPType
	}
}

// Initialize app struct with configuration
func (a *App) Initialize(configuration Configuration) {
	// OVH API
	a.ovh = OVH{}
	a.ovh.Initialize(configuration)
	// Other config
	a.dynDNS.domain = configuration.DynDNS.Domain
	switch strings.ToLower(configuration.DynDNS.Mode) {
	case "dual":
		a.dynDNS.ipTypes = []IPType{IPv4, IPv6}
	case "v4":
		a.dynDNS.ipTypes = []IPType{IPv4}
	case "v6":
		a.dynDNS.ipTypes = []IPType{IPv6}
	}
}

// Run app
func (a *App) Run() {
	Log.Logger.Debug().Msg("Starting DynDNS refresh ...")
	needsRefresh, err := a.needsRefresh()
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Unable to check if DynDNS records need the be renewed.")
		return
	}
	if !needsRefresh {
		Log.Logger.Debug().Msg("DynDNS record refresh not needed.")
		return
	}
	err = a.refresh()
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Error while refreshing DynDNS records.")
	}
	Log.Logger.Debug().Msg("Finished DynDNS refresh.")
}
