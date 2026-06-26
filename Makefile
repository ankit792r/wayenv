# -----------------------------------------------------------------------------
# Project Configuration
# -----------------------------------------------------------------------------

APP_NAME := wayenv
SRC_DIR := source
BIN_DIR := bin

GO := go

VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X main.Version=$(VERSION) \
	-X main.Commit=$(COMMIT) \
	-X main.BuildDate=$(DATE)

BINARY := $(BIN_DIR)/$(APP_NAME)

.DEFAULT_GOAL := help

# -----------------------------------------------------------------------------
# Help
# -----------------------------------------------------------------------------

.PHONY: help
help:
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Build Targets"
	@echo "  build           Build development binary"
	@echo "  release         Build optimized production binary"
	@echo "  build-linux     Build Linux (amd64)"
	@echo "  build-darwin    Build macOS (arm64)"
	@echo "  build-windows   Build Windows (amd64)"
	@echo "  build-all       Build binaries for all supported platforms"
	@echo ""
	@echo "Development"
	@echo "  run             Run the application"
	@echo "  test            Run all tests"
	@echo "  coverage        Run tests with coverage report"
	@echo "  fmt             Format Go source files"
	@echo "  vet             Run go vet"
	@echo "  lint            Run golangci-lint"
	@echo ""
	@echo "Dependencies"
	@echo "  tidy            Clean up go.mod and go.sum"
	@echo "  download        Download Go module dependencies"
	@echo "  install         Install the binary to GOPATH/bin"
	@echo ""
	@echo "Maintenance"
	@echo "  clean           Remove build artifacts"
	@echo "  all             Run fmt, vet, test, and build"
	@echo "  help            Show this help message"
	@echo ""

# -----------------------------------------------------------------------------
# Build
# -----------------------------------------------------------------------------

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && \
	$(GO) build -o ../$(BINARY)

.PHONY: release
release:
	@mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && \
	CGO_ENABLED=0 $(GO) build \
	-trimpath \
	-ldflags "$(LDFLAGS)" \
	-o ../$(BINARY)

.PHONY: run
run:
	cd $(SRC_DIR) && $(GO) run .

# -----------------------------------------------------------------------------
# Testing
# -----------------------------------------------------------------------------

.PHONY: test
test:
	cd $(SRC_DIR) && $(GO) test -v ./...

.PHONY: coverage
coverage:
	cd $(SRC_DIR) && \
	$(GO) test ./... -coverprofile=coverage.out && \
	$(GO) tool cover -html=coverage.out

# -----------------------------------------------------------------------------
# Code Quality
# -----------------------------------------------------------------------------

.PHONY: fmt
fmt:
	cd $(SRC_DIR) && $(GO) fmt ./...

.PHONY: vet
vet:
	cd $(SRC_DIR) && $(GO) vet ./...

.PHONY: lint
lint:
	cd $(SRC_DIR) && golangci-lint run

# -----------------------------------------------------------------------------
# Dependencies
# -----------------------------------------------------------------------------

.PHONY: tidy
tidy:
	cd $(SRC_DIR) && $(GO) mod tidy

.PHONY: download
download:
	cd $(SRC_DIR) && $(GO) mod download

# -----------------------------------------------------------------------------
# Installation
# -----------------------------------------------------------------------------

.PHONY: install
install:
	cd $(SRC_DIR) && $(GO) install

# -----------------------------------------------------------------------------
# Cleanup
# -----------------------------------------------------------------------------

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	cd $(SRC_DIR) && rm -f coverage.out

# -----------------------------------------------------------------------------
# Composite Targets
# -----------------------------------------------------------------------------

.PHONY: build-linux
build-linux:
	@mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -trimpath -ldflags "$(LDFLAGS)" \
	-o ../$(BIN_DIR)/$(APP_NAME)-linux-amd64

.PHONY: build-darwin
build-darwin:
	@mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && \
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
	go build -trimpath -ldflags "$(LDFLAGS)" \
	-o ../$(BIN_DIR)/$(APP_NAME)-darwin-arm64

.PHONY: build-windows
build-windows:
	@mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && \
	GOOS=windows GOARCH=amd64 \
	go build -trimpath -ldflags "$(LDFLAGS)" \
	-o ../$(BIN_DIR)/$(APP_NAME).exe

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: all
all: fmt vet test build
