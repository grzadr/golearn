# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=scalgo
BINARY_UNIX=$(BINARY_NAME)_unix
BIN_DIR=bin

all: test build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v

lint:
	go fmt ./...
	golangci-lint run

test: lint
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

run: build
	./$(BIN_DIR)/$(BINARY_NAME)

# Cross compilation
build-linux: $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/$(BINARY_UNIX) -v

.PHONY: all build lint test clean run build-linux
