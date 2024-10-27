# syntax=docker/dockerfile:1
FROM golang:1.23.2-alpine

WORKDIR /usr/src/youtube2csv

COPY ./go.mod ./go.sum ./

RUN go mod download && go mod verify

COPY ./ ./

# RUN --mount=type=cache,target="/root/.cache/go-build" go build -gcflags "all=-N -l" -o /usr/local/bin/hikariyoga cmd/$HIKARIYOGA_APP_TYPE/main.go
