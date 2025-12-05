FROM golang:1.25.5

ENV CGO_ENABLED=1

RUN apt update && apt upgrade -y
RUN apt install -y keyutils

WORKDIR /usr/src/repo

RUN go install \
   cloudeng.io/go/cmd/gousage@latest \
   cloudeng.io/go/cmd/gomarkdown@latest \
   cloudeng.io/go/cmd/goannotate@latest
RUN go install cloudeng.io/citools/multimod@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
