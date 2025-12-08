# Syntax=docker/dockerfile:1

# Stage 1: Build the binary
FROM golang:1.25 AS builder

WORKDIR /go/src/app

COPY ./*.go .
RUN ls
RUN go mod init cloudeng.io/linux/cmd/docker-entrypoint
RUN go mod download
RUN go mod tidy
RUN ls
RUN go build -o docker-entrypoint .

# Stage 2: Create the final image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates
RUN apk add keyutils

COPY --from=builder /go/src/app/docker-entrypoint /usr/local/bin/docker-entrypoint

# Usage:
# docker build . -t ep-test 
# go run . run --keychain-item=<keys> -- --rm ep-test keyctl list @s
#

ENTRYPOINT ["docker-entrypoint", "entrypoint"]
