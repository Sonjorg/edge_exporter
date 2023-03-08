# Send data to grafana cloud
## A simple setup following this url:
https://grafana.com/docs/grafana-cloud/quickstart/docker-compose-linux/
## Execute following commands from current directory
sudo snap install docker
sudo docker-compose up -d
sudo docker compose ps
## if needed:
sudo docker compose restart
## get ip address of docker container
sudo docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' prometheus

## To see grafana output
curl ip-address-of-prometheus:9090

## or host address in edge
