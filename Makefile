BIN_DIR := "./bin"
BUILD_DIR := "./build"

.PHONY: build
build:
	go build -o $(BIN_DIR)/golinks cmd/golinks/golinks.go

run:
	go run cmd/golinks/golinks.go
