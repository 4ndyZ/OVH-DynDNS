package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Log global logger for the app global
var Log = Logger{}

func main() {
	// Create configuration object
	configuration := Configuration{}
	// Initialize logging
	logFolder := "/var/log/ovh-dyndns"
	_, err := os.Stat(logFolder)
	if os.IsNotExist(err) || os.IsPermission(err) {
		logFolder = "log"
	}
	Log.Initialize(strings.Join([]string{logFolder, "/log.txt"}, ""))

	config := ""
	config1 := "/etc/ovh-dyndns/config.conf"
	config2 := "config/config.conf"
	// Error checking for config files
	_, err1 := os.Stat(config1)
	_, err2 := os.Stat(config2)
	if err1 == nil {
		config = config1
	} else if err2 == nil {
		config = config2
	} else if os.IsNotExist(err1) && os.IsNotExist(err2) {
		Log.Logger.Info().Msg("No configuration file found. Using Commandline parameter.")
	} else if !os.IsNotExist(err1) && os.IsPermission(err1) {
		Log.Logger.Warn().Str("path", config1).Msg("Unable to use configuration file. No permission to access the configuration file.")
	} else if !os.IsNotExist(err2) && os.IsPermission(err2) {
		Log.Logger.Warn().Str("path", config2).Msg("Unable to use configuration file. No permission to access the configuration file.")
	} else if err1 != nil {
		Log.Logger.Warn().Str("error", err1.Error()).Msg("Error while accessing the configuration file.")
	} else if err2 != nil {
		Log.Logger.Warn().Str("error", err2.Error()).Msg("Error while accessing the configuration file.")
	}
	// Try to parse the configuration file if exists
	if config != "" {
		body, err := ioutil.ReadFile(config)
		if err != nil {
			Log.Logger.Warn().Str("error", err.Error()).Msg("Error while reading the configuration file.")
		}
		err = yaml.Unmarshal([]byte(body), &configuration)
		if err != nil {
			Log.Logger.Warn().Str("error", err.Error()).Msg("Error while parsing the configuration file.")
		}
		// Set default configuration parameter
	} else {
		configuration.OVH.Region = "ovh-eu"
		configuration.OVH.ApplicationKey = ""
		configuration.OVH.ApplicationSecretKey = ""
		configuration.OVH.ConsumerKey = ""
		configuration.DynDNS.Domain = "subdomain.example.com"
		configuration.DynDNS.Mode = "v4"
		configuration.DynDNS.Check = "dns"
		configuration.TimeInterval = 60
		configuration.SingleRun = false
		configuration.Logging.Debug = false
	}
	// Commandline flags
	flag.StringVar(&configuration.OVH.Region, "ovh-region", configuration.OVH.Region, "OVH API region [ovh-eu, ovh-us, ovh-ca]")
	flag.StringVar(&configuration.OVH.ApplicationKey, "ovh-ak", configuration.OVH.ApplicationKey, "OVH API application key")
	flag.StringVar(&configuration.OVH.ApplicationSecretKey, "ovh-ask", configuration.OVH.ApplicationSecretKey, "OVH API application secret key")
	flag.StringVar(&configuration.OVH.ConsumerKey, "ovh-ck", configuration.OVH.ConsumerKey, "OVH API consumer key")
	flag.StringVar(&configuration.DynDNS.Domain, "dyndns-domain", configuration.DynDNS.Domain, "Domain for the DynDNS Service")
	flag.StringVar(&configuration.DynDNS.Mode, "dyndns-mode", configuration.DynDNS.Mode, "Mode of the DynDNS service [dual, v4, v6]")
	flag.StringVar(&configuration.DynDNS.Check, "dyndns-check", configuration.DynDNS.Check, "Mode how to check if DynDNS need renewal [api, dns]")
	flag.IntVar(&configuration.TimeInterval, "timeinterval", configuration.TimeInterval, "Time interval in seconds when the service should update the DynDNS records.)")
	flag.BoolVar(&configuration.SingleRun, "singlerun", configuration.SingleRun, "Option to run the microservice only one time and then stop afterwards. Option timeinterval will be ignored!")
	flag.BoolVar(&configuration.Logging.Debug, "debug", configuration.Logging.Debug, "Option to run the microservice in debugging mode")
	flag.Parse()
	// Check if debug log should be enabled
	if configuration.Logging.Debug {
		Log.EnableDebug(true)
	}
	Log.Logger.Info().Msg("Starting ...")
	// Create app worker
	a := App{}
	err = a.Initialize(configuration)
	if err != nil {
		Log.Logger.Error().Str("error", err.Error()).Msg("Unable to start the service.")
		os.Exit(1)
	}
	// Setup signal catching
	sigs := make(chan os.Signal, 1)
	// Catch all signals since not explicitly listing
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	// Method invoked upon seeing signal
	go func() {
		s := <-sigs
		Log.Logger.Info().Str("reason", s.String()).Msg("Stopping the service.")
		Log.Rotate()
		os.Exit(0)
	}()
	// Infinite loop
	for {
		// Run the microservice
		Log.Logger.Info().Msg("Starting refresh ... ")
		a.Run()
		Log.Logger.Info().Msg("Finished refresh.")
		// Check if single run
		if configuration.SingleRun {
			Log.Logger.Info().Msg("Stopping.")
			Log.Rotate()
			os.Exit(0)
		}
		// Wait the provided time to before running again
		d := time.Second * time.Duration(configuration.TimeInterval)
		Log.Logger.Info().Interface("duration", d).Msg("Waiting for the next refresh.")
		time.Sleep(d)
	}

}
