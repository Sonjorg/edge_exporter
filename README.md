# Readme
## Prometheus exporter for Ribbon Communications SBC routers
#### Developed by Sondre Jørgensen in cooperation with Sang Ngoc Nguyen at NTNU: Norwegian University of Science and Technology, sondre2409@gmail.com and 29sangu@gmail.com
### Configuration of the exporter
#### The configuration is implemented in config.yml in the root folder of the source code.
```
---
authtimeout: 3  #all hosts will have max 3 sec timout
hosts:
- hostname: Host1
  ipaddress: 11.111.111.11
  username: Username1
  password: Password1
  routing-database-hours: 24 #For routingentry collector, data is stored in the database for 24 hours for this host.
- hostname: Host2
  ipaddress: 11.111.111.12
  username: Username2
  password: Password2
  routing-database-hours: 24
- hostname: Host3
  ipaddress: 11.111.111.13
  username: Username3
  password: Password3
  routing-database-hours: 24
  exclude:
   - routingentry
   - system
   - diskpartition
   - systemcallstats
   - linecard
   - ethernetport
#Excluding the above collectors for this host
```
- Above you can see the layout of a config.yml file having 3 hosts with dummy data.
- It is required to use a hostname, ipaddress, username and password.
- You can choose which collectors you want to exclude for each host by adding them to the list "exclude" as shown below the last host. The name of the collectors have to match exactly as spelled in this example.
- "Authtimeout" is the maximum chosen time to attempt authentication to a host. Usually it is not reachable if the duration is more than 1-2 second.
- "routing-database-hours" is the duration of which data related to the routingentry collector is stored within the database. Fetching new data through http takes several extra seconds per scrape. Metrics are never stored, only data such as routing tables and their routing entries.
- It is recommended not to use too many hosts per docker instance because of performance issues; a scrape on 2 hosts with no collectors excluded takes around 13 seconds on the first scrape, and around 10 seconds on the following scrapes.

### Deployment running docker
Run:

    sudo docker build -t edge_exporter .

    sudo docker run -p 5123:5123 edge_exporter

Or if you have an external config.yml file:

    sudo docker run -v path/to/your/config.yml:/usr/src/exporter/config.yml sondrjor/edge_exporter

Metrics can be gathered from ```host:5123/metrics```

### Deployment of the SBCexporter on a linux server
**The exporter is developed and tested for the official ubuntu server image found at https://ubuntu.com/download/server.**
- Download golang using the official download page: [install golang](https://go.dev/doc/install), and remember to reboot
- To start the exporter and download all necessary packages, navigate to the SBCexporter directory and run

      go install
### To test go exporters:

    go run .

in the edge_exporter directory, then use

    curl localhost:9100/metrics

in another windows to view live metrics data that can be collected by prometheus

### To test a specific file, for use

    go run main.go
However this will not make use of dependencies from other files

### Installation of Go on HDO's VMs
#### As root folders are not accessible on HDO's VMs we need to install Go in home directory if docker is not utilized
- Download last version of Go to home directory, from Go's official website
- Unzip the file with tar
- Execute the commands:

      export GOPATH=$HOME/go
      export PATH=$PATH:$GOPATH/bin

- If starting Go gives a message that its not yet installed, make a startup script that executes:

      source .bashrc

from home directory

## Grafana and prometheus setup with docker
Choose between grafana local or grafana cloud
### Grafana local
This is a setup with grafana-docker hosted locally, following a similar approach as this tutorial:
https://www.youtube.com/watch?v=9TJx7QTrTyo&t=712s

The config for all docker images used, resides in the docker-compose.yml file

### Deployment of grafana with docker
Use

    docker compose up -d

in either directory ```edge_exporter\Other\Grafana-Prometheus\grafanacloud``` or ```.../grafanalocal```, respectively

## test docker containers:
### get ip address of grafana container

    sudo docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' grafana

    curl ip-address:3000

### Restart all containers if changes are made to docker-compose.yml

    docker-compose up -d --force-recreate

### check status in log files for a container

    sudo docker-compose logs -f container-name

## Collectors and Metrics
List of Collectors, the API endpoints they use and metrics they support:
### System Call Stats Collector
#### Endpoints:
- REST API Method: GET /rest/systemcallstats
#### Metrics:
- Rt_NumCallAttempts - Total number of call attempts system wide since system came up.
- Rt_NumCallSucceeded - Total number of successful calls system wide since system came
up.
- Rt_NumCallFailed - Total number of failed calls system wide since system came up.
- Rt_NumCallCurrentlyUp - Number of currently connected calls system wide.
- Rt_NumCallAbandonedNoTrunk - Number of rejected calls due to no channel available
system wide since system came up.
- Rt_NumCallUnAnswered - Number of unanswered calls system wide since system came
up.
### Diskpartition Collector
#### API endpoints:
- REST API Method: GET /rest/diskpartition - retrieves disk identifier
- REST API Method: GET /rest/diskpartition/{identifier}
#### Metrics:
- Rt_CurrentUsage - Amount of memory used by this partition, expressed as percentage
- Rt_MaximumSize - Specifies the maximum amount of memory, in bytes available in this
partition.
- Rt_MemoryAvailable - Amount of memory in bytes, available for use in the filesystem.
- Rt_MemoryUsed - Amount of memory in bytes, used by the existing files in the
filesystem
- Rt_PartitionName - The name of the disk partition.
- Rt_PartitionType - Identifies the user-friendly physical device holding the partition.
### Ethernetport Collector
#### API endpoints:
- REST API Method: GET /rest/ethernetport - retrieve ethernetport identifier
- REST API Method: GET /rest/ethernetport/{identifier}/historicalstatistics
#### Metrics:
- IfRedundancy - No information found in the Ribbon SBC Edge REST API
Documentation 9.0
- IfRedundantPort - No information found in the Ribbon SBC Edge REST API
Documentation 9.0
- Rt_ifInBroadcastPkts - Displays the number of received broadcast packets on this port.
- Rt_ifInDiscards - Displays the number of discard errors detected on this port.
- Rt_ifInErrors - Displays the number of errors detected on this port.
- Rt_ifInFCSErrors - Displays the number of discard Frame Check Sequence errors
detected on this port.
- Rt_ifInFragmentedPkts - Displays the number of Fragmented Packet errors detected on
this port
- Rt_ifInMulticastPkts - Displays the number of received multicast packets on this port.
- Rt_ifInOctets - Displays the number of received octets on this port.
- Rt_ifInOverSizedPkts - Displays the number of Oversized Packet errors detected on this
port.
- Rt_ifInUcastPkts - Displays the number of received unicast packets on this port.
- Rt_ifInUndersizedPkts - Displays the number of Undersized Packet errors detected on
this port
- Rt_ifInUnknwnProto - Displays the number of Unknown Protocol errors detected on this
port.
- Rt_ifInterfaceIndex - No information found in the Ribbon SBC Edge REST API
Documentation 9.0
- Rt_ifLastChange - The value of sysUpTime at the time the interface entered its current
operational state.
- Rt_ifMtu - The size of the largest packet which can be sent/received on the interface.
- Rt_ifOperatorStatus - The operational status of the interface - 0: IF_OPER_UP or 1:
IF_OPER_DOWN
- Rt_ifOutBroadcastPkts - Displays the number of transmitted broadcast packets on this
port.
- Rt_ifOutDeferredTransmissions - Displays the number of Deferred Transmission errors
detected on this port.
- Rt_ifOutDiscards - Displays the number of discard errors detected on this port.
- Rt_ifOutErrors - Displays the number of errors detected on this port.
- Rt_ifOutLateCollissions - Displays the number of Late Collision errors detected on this
port.
- Rt_ifOutMulticastPkts - Displays the number of transmitted multicast packets on this
port.
- Rt_ifOutOctets - Displays the number of transmitted octets on this port.
- Rt_ifOutUcastPkts - Displays the number of transmitted unicast packets on this port.
- Rt_ifSpeed - An estimate of the interface's current bandwidth in bits per second
- Rt_redundancyRole - When redundancy is configured for "Failover", indicates if it's role
is "Primary" or "Secondary".
- Rt_redundancyState - When redundancy is configured for "Failover", indicates if it's state
is "Online" or "Backup".
### Linecard Collector
#### API endpoints:
REST API Method: GET /rest/linecard/{identifier}
#### Metrics:
- Rt_CardType - The type of hardware module.
- Rt_Location - The hardware module's location within the SBC.
- Rt_ServiceStatus - The service status of the module.
- Rt_Status - Indicates the hardware initialization state for this card.
### Routing Entry Collector
#### API endpoints:
- REST API Method: GET /rest/routingtable
- REST API Method: GET /rest/routingtable/[routingtable]/routingentry
- REST API Method: GET /rest/routingtable/[routingtable]/routingentry/[routingentry]historicalstatistics/1
#### Metrics:
- Rt_RuleUsage - Displays the number of times this call route has been selected for a call.
- Rt_ASR - Displays the Answer-Seizure Ratio for this call route. (ASR is calculated by
dividing the number of call attempts answered by the number of call
attempts.)
- Rt_RoundTripDelay - Displays the average round trip delay for this call route.
- Rt_Jitter - Displays the average jitter for this call route.
- Rt_MOS - Displays the Mean Opinion Score (MOS) for this call route.
- Rt_QualityFailed - Displays if this call route is currently passing or failing the associated
quality metrics. If true then the rule is failing, if false then it is passing.
### System Collector
#### API endpoints:
REST API Method: GET /rest/system/historicalstatistics/1
#### Metrics:
- Rt_CPUUsage - Average percent usage of the CPU.
- Rt_MemoryUsage - Average percent usage of system memory.
- Rt_CPUUptime - The total duration in seconds, that the system CPU has been UP and
running.
- Rt_FDUsage - Number of file descriptors used by the system.
- Rt_CPULoadAverage1m - Average number of processes over the last one minute waiting
to run because CPU is busy.
- Rt_CPULoadAverage5m - Average number of processes over the last five minutes
waiting to run because CPU is busy.
- Rt_CPULoadAverage15m - Average number of processes over the last fifteen minutes
waiting to run because CPU is busy.
- Rt_TmpPartUsage - Percentage of the temporary partition used.
- Rt_LoggingPartUsage - Percentage of the logging partition used. This is applicable only
for the SBC2000.