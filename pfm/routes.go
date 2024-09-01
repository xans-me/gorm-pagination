package pfm

import (
	"net/http"

	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) {
	// Inisialisasi repository dan service
	brimoPFMRepository := NewBrimoPFMRepository()
	brimoPFMService := NewBrimoPFMService(brimoPFMRepository)
	brimoPFMController := NewBrimoPFMController(brimoPFMService, db)

	// Define route and handler
	http.HandleFunc("/api/brimo_pfm", brimoPFMController.GetBrimoPFMHandler)
}
