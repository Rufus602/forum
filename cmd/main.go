package main

import (
	"log"
	"testForum/internal/server"
)

func main() {
	err := server.Server()
	if err != nil {
		log.Fatal(err)
	}
}
