# syntax=docker/dockerfile:1
FROM golang:alpine

ENV WORKDIR=/usr/src/youtube2csv

WORKDIR ${WORKDIR}

COPY ./go.mod ./go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /usr/local/bin/youtube2csv cmd/main.go

CMD ["youtube2csv"]
