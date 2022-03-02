package main

// IPType used as enum for IP address type (IPv4 or IPv6)
type IPType string

// IP Type enum constants
const (
	IPv4 = "A"
	IPv6 = "AAAA"
)

func (ipType IPType) String() string {
	return string(ipType)
}

// Configuration struct
type Configuration struct {
	OVH struct {
		Region               string `yaml:"region"`
		ApplicationKey       string `yaml:"applicationKey"`
		ApplicationSecretKey string `yaml:"applicationSecretKey"`
		ConsumerKey          string `yaml:"consumerKey"`
	} `yaml:"ovh"`
	DynDNS struct {
		Domain string `yaml:"domain"`
		Mode   string `yaml:"mode"`
	} `yaml:"dyndns"`
	TimeInterval int  `yaml:"timeinterval-to-pull"`
	SingleRun    bool `yaml:"single-run"`
	Logging      struct {
		Debug bool `yaml:"debug"`
	} `yaml:"logging"`
}
