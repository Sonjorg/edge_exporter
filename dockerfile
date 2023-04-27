FROM golang:1.20 as build

WORKDIR /usr/src/exporter

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/exporter ./

CMD ["exporter"]

EXPOSE 5123