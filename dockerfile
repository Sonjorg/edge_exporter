FROM golang:1.20

WORKDIR /var/exporter

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#RUN mkdir exporter
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN chmod 777 /var
RUN chmod 777 /var/exporter

RUN go build -o main.go

CMD ["/var/exporter"]

EXPOSE 5123