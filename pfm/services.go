package pfm

import (
	"net/http"
	"strconv"
	"test-pagination-pg-go/pagination"
)

func GetPaginatedTransactions(r *http.Request) (interface{}, error) {
	db := GetDB()

	// Ambil query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}

	sort := r.URL.Query()["sort"]
	dateStart := r.URL.Query().Get("dateStart")
	dateEnd := r.URL.Query().Get("dateEnd")
	accountNumber := r.URL.Query().Get("account_number")
	trxAmountStr := r.URL.Query().Get("trx_amount")

	// Konversi trxAmount ke float64
	var trxAmount float64
	if trxAmountStr != "" {
		trxAmount, _ = strconv.ParseFloat(trxAmountStr, 64)
	}

	// Buat FilterManager untuk mengelola filter
	filterManager := pagination.FilterManager{}

	// Tambahkan filter untuk trx_date
	if dateStart != "" && dateEnd != "" {
		filterManager.AddAndFilter(pagination.DateRangeFilter{
			Field:     "trx_date",
			StartDate: dateStart,
			EndDate:   dateEnd,
		})
	}

	// Tambahkan filter untuk account_number
	if accountNumber != "" {
		filterManager.AddAndFilter(pagination.ComparisonFilter{
			Field:    "account_number",
			Operator: "=",
			Value:    accountNumber,
		})
	}

	// Tambahkan filter untuk trx_amount
	if trxAmount > 0 {
		filterManager.AddAndFilter(pagination.ComparisonFilter{
			Field:    "trx_amount",
			Operator: ">=",
			Value:    trxAmount,
		})
	}

	// Tambahkan filter OR untuk trx_type (pengeluaran OR pemasukan)
	filterManager.AddOrFilter(pagination.ComparisonFilter{
		Field:    "trx_type",
		Operator: "=",
		Value:    "pengeluaran",
	})

	filterManager.AddOrFilter(pagination.ComparisonFilter{
		Field:    "trx_type",
		Operator: "=",
		Value:    "pemasukan",
	})

	// Terapkan filter ke query
	query := filterManager.Apply(db.Model(&BrimoPFM{})).Debug()

	// Terapkan sorting dan pagination
	paginator := pagination.NewPaginator(
		query,
		pagination.WithPage(page),
		pagination.WithPageSize(pageSize),
		pagination.WithSort(sort...),
	)

	var transactions []BrimoPFM
	result, err := paginator.Paginate(&transactions)
	if err != nil {
		return nil, err
	}

	return result, nil
}
