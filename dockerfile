FROM golang:1.20 as build

WORKDIR /usr/src/exporter

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ARG authtimeout
ARG hostname
ARG ipaddress
ARG username
ARG password
ARG exclude
ARG routing_database_hours

ENV authtimeout $authtimeout
ENV hostname $hostname
ENV ipaddress $ipaddress
ENV username $username
ENV password $password
ENV exclude $exclude
ENV routing_database_hours $routing_database_hours

COPY . .

RUN go build -v -o /usr/local/bin/exporter ./

CMD ["exporter"]

EXPOSE 5123