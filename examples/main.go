package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"test-pagination-pg-go/examples/transaction"
)

func main() {
	r := mux.NewRouter()

	// Register routes
	transaction.RegisterRoutes(r)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
