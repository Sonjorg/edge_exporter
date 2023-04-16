FROM golang:1.20

WORKDIR /var/exporter

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#RUN mkdir exporter
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN chmod 777 /usr

RUN go build -o main.go

#RUN chmod -x /
CMD ["/var/exporter"]

EXPOSE 5123