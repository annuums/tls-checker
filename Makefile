GOPATH := $(shell go env GOPATH)
APP ?= tls-checker
CMD_DIR := ./cmd/$(APP)

.PHONY: build build-linux.amd64 build-linux.arm64 build-macos test

build-macos:
	mkdir -p deploy/bin
	CGO_ENABLED=0 go build -a -installsuffix cgo \
		-ldflags="-w -s" -o deploy/bin/main.macos $(CMD_DIR)

build-linux.amd64:
	mkdir -p deploy/bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
		-ldflags="-w -s" -o deploy/bin/main.linux.amd64 $(CMD_DIR)

build-linux.arm64:
	mkdir -p deploy/bin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo \
		-ldflags="-w -s" -o deploy/bin/main.linux.arm64 $(CMD_DIR)

test:
	go test -v ./... -short
