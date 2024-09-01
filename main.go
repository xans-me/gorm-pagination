package main

import (
	"log"
	"net/http"
	"test-pagination-pg-go/pfm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=pfm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	} // Aktifkan debug mode
	db = db.Debug()

	// Setup routes
	pfm.SetupRoutes(db)

	// Start server
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
