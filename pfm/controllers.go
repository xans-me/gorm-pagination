package pfm

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type BrimoPFMController struct {
	service BrimoPFMService
	db      *gorm.DB
}

func NewBrimoPFMController(service BrimoPFMService, db *gorm.DB) *BrimoPFMController {
	return &BrimoPFMController{service: service, db: db}
}

func (c *BrimoPFMController) GetBrimoPFMHandler(w http.ResponseWriter, r *http.Request) {
	// Extract query params for pagination, sorting, filtering, etc.
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	sortField := r.URL.Query().Get("sort_field")
	if sortField == "" {
		sortField = "trx_date"
	}

	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	trxType := r.URL.Query().Get("trx_type")
	if trxType == "" {
		trxType = "pemasukan"
	}

	cif := r.URL.Query().Get("cif")
	accountNumber := r.URL.Query().Get("account_number")
	groupBy := r.URL.Query().Get("group_by")
	prefetch := r.URL.Query().Get("prefetch") == "true"

	filters := map[string]interface{}{
		"trx_type": trxType,
	}

	if cif != "" {
		filters["cif"] = cif
	}

	if accountNumber != "" {
		filters["account_number"] = accountNumber
	}

	response, err := c.service.GetBrimoPFM(c.db, page, pageSize, sortField, sortOrder, groupBy, filters, prefetch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
