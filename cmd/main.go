package main

import (
	"fmt"
	"log"
	"os"

	"github.com/samthehai/chat/internal/application/config/wire"
)

func main() {
	server, cleaner, err := wire.InitializeServer()
	if err != nil {
		log.Fatalf("failed to create server: %+v", err)
	}

	if err := server.Serve(); err != nil {
		cleaner()
		fmt.Fprintf(os.Stderr, "failed to run server: %+v", err)
	}
}
