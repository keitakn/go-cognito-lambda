FROM golang:1.15-alpine3.12

LABEL maintainer="https://github.com/keitakn"

WORKDIR /go/app

COPY . .

ENV GO111MODULE=off

ARG GOLANGCI_LINT_VERSION=v1.34.0

RUN set -eux && \
  apk update && \
  apk add --no-cache git curl make && \
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION} && \
  go get golang.org/x/tools/cmd/goimports

ENV GO111MODULE on
ENV CGO_ENABLED 0
