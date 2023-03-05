# Readme
## SBC Exporter written in golang
#### The prometheus golang exporter resides in the SBCexporter folder

### Deployment of the SBCexporter
- Download golang using the official download page and remember to reboot
- To start the exporter and download all necessary packages, navigate to the SBCexporter directory and run
``` go install ```
### To test go exporters:
``` go run . ``` in the SBCexporter directory, then use ```curl localhost:9100/metrics``` in another windows to view live metrics data that can be collected by prometheus
### To test a specific file, for use
``` go run main.go ``` However this will not make use of dependencies from other files


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
