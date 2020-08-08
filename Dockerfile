#!/bin/bash

FROM golang:latest as builder

WORKDIR /go/src/TakkenGo
ENV GO111MODULE on
COPY ./go.mod  ./
RUN go mod download
COPY . .

WORKDIR ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main

FROM alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/TakkenGo/cmd/main /main
RUN chmod +x /main

ENV PORT 8080

CMD ["/main"]
