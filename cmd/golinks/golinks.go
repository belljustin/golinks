package main

import (
	"log"

	"github.com/belljustin/golinks"
	"github.com/belljustin/golinks/internal/storage/memory"
)

func main() {
	log.Println("Starting golinks server...")

	storage := memory.NewStorage()
	server := golinks.NewServer(storage, "./web")
	log.Fatal(server.Serve())
}
