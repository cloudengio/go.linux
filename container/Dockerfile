FROM golang:1.25.5-alpine3.23

RUN apk add keyutils

WORKDIR /usr/src/repo

RUN go install cloudeng.io/citools/multimod@latest
RUN go install cloudeng.io/go/cmd/gomarkdown@latest
RUN go install golang.org/x/vuln/cmd/govulncheck@latest
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
