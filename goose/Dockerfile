FROM golang:1.25-alpine

RUN apk add --no-cache git make bash

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ENTRYPOINT ["goose"]
