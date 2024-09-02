package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=pfm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	} // Aktifkan debug mode
	db = db.Debug()

	// Start server
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
