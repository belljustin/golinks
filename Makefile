BIN_DIR := "./bin"
BUILD_DIR := "./build"

.PHONY: build
build:
	go build -o $(BIN_DIR)/golinks cmd/golinks/golinks.go

docker-build:
	docker build -t golinks:latest -t golinks:stable -f $(BUILD_DIR)/docker/Dockerfile .