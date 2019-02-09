package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Peltoche/avro-gateway/registry"
	"github.com/Peltoche/avro-gateway/schema"
	"github.com/gorilla/mux"
)

const addr = ":8080"

func main() {
	router := mux.NewRouter()

	schemaRegistryURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		panic(err)
	}

	// Clients.
	registry := registry.NewClient(schemaRegistryURL)

	// Schema.
	schemaUsecase := schema.NewUsecase(registry)
	schemaHandler := schema.NewHTTPHandler(schemaUsecase)
	schemaHandler.RegisterRoutes(router)

	log.Printf("start listening on %s", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
