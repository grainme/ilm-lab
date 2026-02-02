package main

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Starting server on port %s", port)
}
