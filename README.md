# Readme
## SBC Exporter written in golang
**The prometheus golang exporter resides in the SBCexporter folder**

### Deployment of the exporter running docker
- Run:
``` sudo docker build -t exporter .```
``` sudo docker run -p 5123:5123 exporter ```
- Metrics can be gathered from ```host:5123/metrics```

### Configuration of the exporter
##### The configuration is implemented in config.yml in the root folder of the source code.
- Below you can see the layout of a config.yml file having 3 hosts. 
- It is required to use a hostname, ipaddress, username and password. 
- You can choose which collectors you want to exclude for each host by adding them to the list "exclude" as shown below the last host.  
- "Authtimeout" is the maximum chosen time to attempt authentication to a host. Usually it is not reachable if the duration is more than 1-2 second. 
- "routing-database-hours" is the duration of which data related to the routingentry collector is stored within the database. Fetching new data through http takes several extra seconds per scrape. Metrics are never stored, only data such as routing tables and their routing entries.
- It is recommended not to use too many hosts per docker instance because of performance issues; a scrape on 2 hosts with no collectors excluded takes around 13 seconds on the first scrape, and around 9 seconds on the following scrapes.
```
---
authtimeout: 3  #all hosts will have max 3 sec timout
hosts:
- hostname: Host1
  ipaddress: 11.111.111.11
  username: Username1
  password: Password1
  routing-database-hours: 24 #For routingentry collector, data is stored in the database for 24 hours for this host, NB: can not be 0
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

### Deployment of the SBCexporter on a linux server
**The exporter is developed and tested for the official ubuntu server image found at https://ubuntu.com/download/server.**
- Download golang using the official download page: [install golang](https://go.dev/doc/install), and remember to reboot
- To start the exporter and download all necessary packages, navigate to the SBCexporter directory and run
``` go install ```
### To test go exporters:
``` go run . ``` in the SBCexporter directory, then use ```curl localhost:9100/metrics``` in another windows to view live metrics data that can be collected by prometheus
### To test a specific file, for use
``` go run main.go ``` However this will not make use of dependencies from other files

### Installation of Go on HDO's VMs
#### As root folders are not accessible on HDO's VMs we need to install Go in home directory if docker is not utilized
- Download last version of Go to home directory, from Go's official website
- Unzip the file with tar
- Execute the commands:
``` export GOPATH=$HOME/go ```
``` export PATH=$PATH:$GOPATH/bin ```
- If starting Go gives a message that its not yet installed, make a startup script that executes:
``` source .bashrc ``` from home directory

## Grafana and prometheus setup with docker

This is a setup with grafana-docker hosted locally, following a similar approach as this tutorial:
https://www.youtube.com/watch?v=9TJx7QTrTyo&t=712s

The config for all docker images used, resides in the docker-compose.yml file

### Deployment of docker images
Use
``` docker compose up -d ```
in the current dir


## test docker:
### get ip address of grafana container
``` sudo docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' grafana ```

curl ip-address:3000

### Restart all containers if changes are made to docker-compose.yml
``` docker-compose up -d --force-recreate ```

### check status in log files for a container
```sudo docker-compose logs -f container-name ```
