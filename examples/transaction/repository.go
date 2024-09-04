package transaction

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	// Setup DB connection (using PostgresSQL)
	db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	db.Debug()
}

// GetDB returns the database instance.
func GetDB() *gorm.DB {
	return db
}
