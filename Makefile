# Makefile for managing the build process of a gRPC application in Go

# Variables
PROTO_FILE := ./rpc/daemon.proto
PROTO_GEN_GO_DIR := ./pb
GO_CMD := go
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


build: build_daemon build_client

# Build the server
build_daemon:
	$(GO_CMD) build -o sshtmd ./daemon

# Build the client
build_client:
	$(GO_CMD) build -o sshtm ./client

# Clean up generated files and binaries
clean:
	rm -f ./sshtmd ./sshtm
	rm -rf $(PROTO_GEN_GO_DIR)/*

# Help
help:
	@echo "Usage:"
	@echo "  make gen_proto     Generate Go code from .proto file"
	@echo "  make build_daemon  Build the gRPC daemon"
	@echo "  make build_client  Build the gRPC client"
	@echo "  make clean         Clean up generated files and binaries"
