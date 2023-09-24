binary_name=main

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GRPC_MODULE ?= twofa
BUILD_FLAGS ?= -v

WINDOWS_EXEC := main.exe
LINUX_EXEC := main
MAC_EXEC := main

all: windows linux mac

windows:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ./bin/$(WINDOWS_EXEC) ./cmd/api/main.go
	
linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ./bin/$(LINUX_EXEC) ./cmd/api/main.go
	
mac:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ./bin/$(MAC_EXEC) ./cmd/api/main.go

run-windows:
	./bin/${WINDOWS_EXEC}

run-linux:
	./bin/${LINUX_EXEC}

run-mac:
	./bin/${MAC_EXEC}

clean:
	rm -rf ./bin/*

dev:
	go run ./cmd/server/main.go

fmt:
	gofmt -w .

Test:
	go clean -testcache
	go test ./... -race -v

proto:
	protoc -I=internal/proto \
	--go_out=internal/proto \
	--go-grpc_out=internal/proto \
	internal/proto/${GRPC_MODULE}.proto