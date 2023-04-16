FROM golang:1.20 as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
#RUN chmod 777 /var
#RUN chmod 777 /var/exporter

#COPY --from=build /app /var/exporter
RUN go build -v -o /usr/local/bin/app ./
#RUN go build -o main.go

CMD ["app"]

EXPOSE 5123