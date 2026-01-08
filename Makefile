VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/naxodev/github-switch/cmd.Version=$(VERSION)"

.PHONY: build install clean test

build:
	go build $(LDFLAGS) -o github-switch .

install:
	go install $(LDFLAGS) .

clean:
	rm -f github-switch

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

.DEFAULT_GOAL := build
