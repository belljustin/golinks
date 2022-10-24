package main

import (
	"fmt"
	"log"
	"os"

	"github.com/belljustin/golinks/internal/golinks"
	_ "github.com/belljustin/golinks/internal/storage/dynamodb"
	_ "github.com/belljustin/golinks/internal/storage/memory"
)

func main() {
	if len(os.Args) == 1 {
		serve()
	}

	switch os.Args[1] {
	case "migrate":
		migrate()
	default:
		panic(fmt.Sprintf("Unrecognized command %s", os.Args[1]))
	}
}

func migrate() {
	storage := golinks.NewStorage()
	if err := storage.Migrate(); err != nil {
		panic(err)
	}
}

func serve() {
	server := golinks.NewServer()
	log.Fatal(server.Serve())
}
