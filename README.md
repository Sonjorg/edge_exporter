## Simple setup with grafana hosted locally

## Following a similar approach as this tutorial:
https://www.youtube.com/watch?v=9TJx7QTrTyo&t=712s

Use
``` docker compose up -d ```
in the current dir

## test to see if it works:
### get ip address of grafana container
``` sudo docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' grafana ```

curl ip-address:3000

### Restart all containers if changes are made to docker-compose.yml
``` docker-compose up -d --force-recreate ```

### check status in log files for a container
```sudo docker-compose logs -f container-name ```
### Download golang using the official download page and remember to reboot