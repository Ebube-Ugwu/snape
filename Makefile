GOCACHE ?= /tmp/snape-go-build
BINARY ?= dist/snape
MAIN := ./cmd

.PHONY: all build run serve test fmt tidy clean

all: fmt test build

build:
	mkdir -p dist
	GOCACHE=$(GOCACHE) go build -o $(BINARY) $(MAIN)

run:
	GOCACHE=$(GOCACHE) go run $(MAIN) list

serve:
	GOCACHE=$(GOCACHE) go run $(MAIN) serve --addr :7777

test:
	GOCACHE=$(GOCACHE) go test ./...

fmt:
	gofmt -w cmd internal

tidy:
	go mod tidy

clean:
	rm -rf dist
