FROM golang:1.16-alpine

WORKDIR /SBCexporter

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY /SBCexporter ./exporter

RUN go build -o ./exporter

COPY --from=build /exporter /exporter

EXPOSE 9100
#ENTRYPOINT ["/exporter"]