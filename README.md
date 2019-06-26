# What?
A continuously expanding set of utilities for interacting with Palo Alto Networks devices.

# Why?
This is a place to put handy stuff that doesn't already have a home, or for prototyping new functions.

It's useful to have a set of functions that are OS indepdent and require no installer. 

# How?
## Install
Panutil is bundled as a single executable for each platform on the [releases page](https://github.com/adambaumeister/panfw-util/releases).

Simply download the version for your OS and proceed to the setup section.

## Setup
Before getting started with panutil, you must provide it a login to the device you want to target 
as well as it's address.

The options are;

- A YAML configuration file [like this one](.panutil.yml)
- Environment variables
    ```bash
    PANUTIL_USER
    PANUTIL_PASSWORD
    PANUTIL_HOSTNAME
    ```
- Command line flags
    ```bash
    panutil.exe --username blah --password blah --hostname blah
    ```

HOSTNAME is an FQDN or IP address, with an optional port such as _test.corp.local:443_ or _192.168.1.1_
 
## Functions
### Help
_help_ allows you to navigate around panutil's CLI. 
```bash
panutil.exe -h
```

### Load/Import
_load_ allows you to load and, optionally, commit an XML configuration file on local disk to a PANOS device.

This is useful for standing up environments from scratch as part of a greater automation pipeline.
```bash
panutil load running-config.xml
```

### Add
_add_ lets you bulk add supported object types. 

Currently, only CSV is supported as the source for add.

```bash
# Add an address object
panutil add address,test,1.1.1.1
# Add an object to an address group and create the group if not already existing
panutil add addressgroup,testgroup,test
# Add a service
panutil add service,test-service,tcp,8080
```

### Testpcap
_testpcap_ retrieves flows from a saved packet capture in pcap format and runs the PANOS test functionality based on the flows seen.

This is useful for testing policy matches against real traffic. Optionally, source and destination zone can be specified to increase the accuracy of the tests.

```bash
panutil testpcap capture.pcap 
```

# Development
Get the code and dependencies. This will put all the code in %GOPATH%/github.com/adambaumeister/panfw-util
```
go get github.com/adambaumeister/panfw-util
go get ./...
```

Make your changes and submit a PR on Github.



