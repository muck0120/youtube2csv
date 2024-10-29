# syntax=docker/dockerfile:1
FROM golang:alpine

WORKDIR /usr/src/youtube2csv

# COPY ./go.mod ./go.sum ./
COPY ./go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /usr/local/bin/youtube2csv cmd/main.go

CMD ["youtube2csv"]
