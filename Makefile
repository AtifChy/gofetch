# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=gofetch

# Source directory
SRC_DIR=./cmd/gofetch

# Build directory
BUILD_DIR=bin

.PHONY: all build clean test coverage deps run help install uninstall

## all: Default target
all: clean deps test build

## build: Build the binary
build:
	$(GOBUILD) -o $(BUILD_DIR)/ -v $(SRC_DIR)

## clean: Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

## test: Run tests
test:
	$(GOTEST) -v ./...

## deps: Download and tidy dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## run: Build and run the binary
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

## install: Install binary to GOPATH/bin
install:
	$(GOMOD) install $(SRC_DIR)

## uninstall: Remove binary from GOPATH/bin
uninstall:
	rm -f $(GOPATH)/$(BUILD_DIR)/$(BINARY_NAME)

## help: Show help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /' | sed 's/: / - /'
