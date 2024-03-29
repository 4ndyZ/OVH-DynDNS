package main

import (
	"errors"
	"strings"
)

// App struct to hold refs
type App struct {
	ovh    OVH
	dynDNS struct {
		domain      string
		ipTypes     []IPType
		checkMethod CheckMethod
	}
}

// Initialize app struct with configuration
func (a *App) Initialize(configuration Configuration) error {
	// OVH API
	a.ovh = OVH{}
	err := a.ovh.Initialize(configuration)
	if err != nil {
		return errors.New("error while loading OVH configuration")
	}
	// Other config
	a.dynDNS.domain = configuration.DynDNS.Domain
	switch strings.ToLower(configuration.DynDNS.Mode) {
	case "dual":
		a.dynDNS.ipTypes = []IPType{IPv4, IPv6}
	case "ipv4":
		a.dynDNS.ipTypes = []IPType{IPv4}
	case "ipv6":
		a.dynDNS.ipTypes = []IPType{IPv6}
	default:
		Log.Logger.Warn().Str("mode", strings.ToLower(configuration.DynDNS.Mode)).Msg("Invalid DynDNS service mode. Defaulting to \"ipv4\".")
		a.dynDNS.ipTypes = []IPType{IPv4}
	}
	Log.Logger.Debug().Str("mode", strings.ToLower(configuration.DynDNS.Mode)).Msg("Running mode.")
	switch strings.ToLower(configuration.DynDNS.Check) {
	case "dns":
		a.dynDNS.checkMethod = DNS
	case "api":
		a.dynDNS.checkMethod = API
	default:
		Log.Logger.Warn().Str("check", strings.ToLower(configuration.DynDNS.Check)).Msg("Invalid DynDNS check method. Defaulting to \"DNS\".")
	}
	Log.Logger.Debug().Str("check", strings.ToLower(configuration.DynDNS.Check)).Msg("Checking mode.")
	return nil
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
