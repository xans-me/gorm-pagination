package main

import (
	"github.com/gorilla/mux"
	"github.com/xans-me/gorm-pagination/examples/transaction"
	"log"
	"net/http"
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
