# Makefile for managing the build process of a gRPC application in Go

# Variables
PROTO_FILE := ./rpc/daemon.proto
PROTO_GEN_GO_DIR := ./pb
GO_CMD := go
GOPROXY ?= https://proxy.golang.org,direct
export GOPROXY
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null | sed "s/^v//")
LDFLAGS := $(if $(VERSION),-ldflags "-X github.com/besrabasant/ssh-tunnel-manager/config.AppVersion=$(VERSION)")
PROTOC_CMD := protoc
PROTO_PATH := . # Adjust this path to where your .proto files are located
GO_OUT := .
GRPC_OUT := .

# Ensure that $(PROTO_GEN_GO_DIR) exists
$(shell mkdir -p $(PROTO_GEN_GO_DIR))

# Default target
all: gen_proto build

# Generate Go code from .proto file
gen_proto:
	$(PROTOC_CMD) --go_out=$(GO_OUT) --go_opt=paths=source_relative --go-grpc_out=$(GRPC_OUT) --go-grpc_opt=paths=source_relative $(PROTO_FILE)


download_deps:
	$(GO_CMD) mod download

build: build_daemon build_client

# Build the server (current OS/arch)
build_daemon:
	$(GO_CMD) build $(LDFLAGS) -o sshtmd ./daemon

# Build the client (current OS/arch)
build_client:
	$(GO_CMD) build $(LDFLAGS) -o sshtm ./client

# Cross-platform builds (macOS + Linux)
build_linux:
	GOOS=linux GOARCH=amd64 $(GO_CMD) build $(LDFLAGS) -o sshtmd-linux-amd64 ./daemon
	GOOS=linux GOARCH=amd64 $(GO_CMD) build $(LDFLAGS) -o sshtm-linux-amd64 ./client

build_macos:
	GOOS=darwin GOARCH=amd64 $(GO_CMD) build $(LDFLAGS) -o sshtmd-darwin-amd64 ./daemon
	GOOS=darwin GOARCH=amd64 $(GO_CMD) build $(LDFLAGS) -o sshtm-darwin-amd64 ./client
	GOOS=darwin GOARCH=arm64 $(GO_CMD) build $(LDFLAGS) -o sshtmd-darwin-arm64 ./daemon
	GOOS=darwin GOARCH=arm64 $(GO_CMD) build $(LDFLAGS) -o sshtm-darwin-arm64 ./client

build_all: build_linux build_macos

# Clean up generated files and binaries
clean:
	rm -f ./sshtmd ./sshtm
	rm -rf $(PROTO_GEN_GO_DIR)/*

# Help
help:
	@echo "Usage:"
	@echo "  make gen_proto     Generate Go code from .proto file"
	@echo "  make download_deps  Download Go module dependencies"
	@echo "  make build_daemon  Build the gRPC daemon"
	@echo "  make build_client  Build the gRPC client"
	@echo "  make build_linux   Build Linux binaries (amd64)"
	@echo "  make build_macos   Build macOS binaries (amd64 + arm64)"
	@echo "  make build_all     Build Linux + macOS binaries"
	@echo "  make VERSION=1.2.3 build_all  Build with an explicit app version"
	@echo "  make clean         Clean up generated files and binaries"
