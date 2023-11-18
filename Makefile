# Include the .env file
include .env

# Variables
BINARY_NAME := main
DOCKER_PROJECT := nestpass-auth
BUILD_FLAGS ?= -v
GRPC_MODULE ?= 

OUTPUT_EXEC := ./bin/\$(BINARY_NAME)

# Build Target
build:
	@echo "Building..."
	@rm -rf ./bin/*
	@rm -rf ./bin/.gitkeep
	@go build \$(BUILD_FLAGS) -o \$(OUTPUT_EXEC) ./cmd/server/main.go

docker-build:
	@echo "Removing previous docker image if it exists..."
	@docker rmi -f \$(DOCKER_PROJECT) || true
	@echo "Building docker image..."
	@docker build -t \$(DOCKER_PROJECT) .

# Run Target
run:
	@echo "Running binary..."
	@\$(OUTPUT_EXEC)

docker-run:
	@echo "Running docker image..."
	@docker run --name \$(DOCKER_PROJECT) -p \$(PORT):\$(PORT) --env-file .env \$(DOCKER_PROJECT)

# Development and Utility Targets
clean:
	@echo "Cleaning up binaries..."
	@rm -rf ./bin/*

dev:
	@echo "Starting the dev server..."
	@go run ./cmd/server/main.go

fmt:
	@echo "Formatting the code..."
	@gofmt -w .

run-tests:
	@echo "Running tests..."
	@go clean -testcache
	@go test ./... -race -v

proto:
	@echo "Generating gRPC code..."
	@protoc -I=internal/proto \
	--go_out=internal/proto \
	--go-grpc_out=internal/proto \
	internal/proto/\$(GRPC_MODULE).proto