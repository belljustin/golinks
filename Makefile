BIN_DIR := "./bin"

run: build
	$(BIN_DIR)/golinks

.PHONY: build
build:
	go build -o $(BIN_DIR)/golinks cmd/golinks/golinks.go

clean:
	rm -rf $(BIN_DIR)