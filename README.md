# DynDNS Client for OVH
This microservice uses the official [OVH API](https://api.ovh.com/) and [their Go module](https://github.com/ovh/go-ovh/) to update a DynDNS Record.

## Function
The microservice supports IPv4 (A record) and IPv6 (AAAA record). The microservice obtains you public IPv4 and IPv6 request by connection to a webserver and then creates or updates the DNS Zone via the OVH API.

## Notice of Non-Affiliation and Disclaimer
This project is not affiliated, associated, authorized, endorsed by, or in any way officially connected with OVH SAS, or any of its subsidiaries or its affiliates.

## Prerequisites
It is recommended to have a server where you can deploy the microservice, but it is also possible to start the microservice manually on a local machine.

## Installation and configuration
Download the prebuilt binary packages from the [release page](https://github.com/4ndyZ/OVH-DynDNS/releases) and install them on your server.

### Installation
#### Linux
###### DEB Package
If you are running a Debian-based Linux Distribution choose the `.deb` Package for your operating system architecture and download it. You are able to use curl to download the package.

Now you are able to install the package using APT.
`sudo apt install ./ovh-dyndns-vX.X-.linux.XXXX.deb`

After installing the package configure the microservice. The configuration file is located under `/etc/ovh-dyndns/config.yml`.

At this point you are able to enable the Systemd service using `systemctl`.
`sudo systemctl enable ovh-dyndns`

And start the service also using `systemctl`.
`sudo systemctl start ovh-dyndns`

###### RPM Package
When running a RHEL-based Linux Distribution choose the `.rpm` package for your operating system architecture and download it.

Now you are able to install the package.
`sudo dnf install ./ovh-dyndns-vX.X-.linux.XXXX.rpm`

After installing the package configure the microservice. The configuration file is located under `/etc/ovh-dyndns/config.yml`.

Now you are able to enable the Systemd service using `systemctl`.
`sudo systemctl enable ovh-dyndnss`

And start the service also using `systemctl`.
`sudo systemctl start ovh-dyndns`

#### Windows/Other
If you plan to run the microservice on Windows or another OS the whole process is a bit more complicated because there is no installation package available only prebuilt binaries.

Download the prebuilt binary for your operating system.

Extract the prebuilt binary and change the configuration file located under `config/config.conf`.

After successful changing the configuration file you are able to run the prebuilt binary.

### Configuration
The microservice tries to access the configuration file located under `/etc/ovh-dyndns/config.conf`. If the configuration file is not accessible or found the service will fall back to the local file located under `config/config.conf`.
To get the required application key, application secret key and consumer key visit the [OVH API documentation](https://docs.ovh.com/gb/en/api/first-steps-with-ovh-api/#create-your-app-keys).

When creating the API key you need the following permission:

```
GET: /domain/zone/*
PUT: /domain/zone/*
POST: /domain/zone/*
```

Or if you want to restrict access to one domain zone:

```
GET: /domain/zone/{domain-zone}/*
PUT: /domain/zone/{domain-zone}/*
POST: /domain/zone/{domain-zone}/*
```

### Logging
The microservice while try to put the log file in the `/var/log/ovh-dyndns` folder. If the service is not able to access or find that folder, the logging file gets created in the local folder `logs`.

If you want to enable debug messages please change the configuration file  or run the microservice with the commandline parameter `-debug`.

## Usage
```
Usage:
  -debug
        Option to run the microservice in debugging mode
  -dyndns-check string
        Mode how to check if DynDNS need renewal [api, dns] (default "api")
  -dyndns-domain string
        Domain for the DynDNS Service (default "subdomain.example.com")
  -dyndns-mode string
        Mode of the DynDNS service [dual, ipv4, ipv6] (default "ipv4")
  -network-interface string
        Network interface to use for the DynDNS service. If empty the auto selected network interface will be used.
  -ovh-ak string
        OVH API application key
  -ovh-ask string
        OVH API application secret key
  -ovh-ck string
        OVH API consumer key
  -ovh-region string
        OVH API region [ovh-eu, ovh-us, ovh-ca] (default "ovh-eu")
  -singlerun
        Option to run the microservice only one time and then stop afterwards. Option timeinterval will be ignored!
  -timeinterval int
        Time interval in seconds when the service should update the DynDNS records.) (default 120)
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[GPL-3.0](https://github.com/4ndyZ/OVH-DynDNS/blob/main/COPYING)
