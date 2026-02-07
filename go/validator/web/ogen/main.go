package main

import (
	"log"
	"net/http"

	"github.com/buzztaiki/sandbox/go/validator/web/ogen/petstore"
)

func main() {
	// Create service instance.
	service := &petsService{
		pets: map[int64]petstore.Pet{},
	}
	// Create generated server.
	srv, err := petstore.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
