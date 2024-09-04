package pfm

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"test-pagination-pg-go/pagination"
)

func GetPaginatedTransactions(r *http.Request) (interface{}, error) {
	db := GetDB()

	// Parse query parameters
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
	search := r.URL.Query().Get("search")

	// Convert transaction amount to float64
	var trxAmount float64
	if trxAmountStr != "" {
		trxAmount, _ = strconv.ParseFloat(trxAmountStr, 64)
	}

	// Initialize base query
	query := db.Model(&TransactionData{})

	// Initialize FilterManager
	filterManager := pagination.FilterManager{}

	// Add filters to FilterManager
	addDateRangeFilter(&filterManager, dateStart, dateEnd)
	addTransactionAmountFilter(&filterManager, trxAmount)

	// Apply filters from FilterManager
	query = filterManager.Apply(query)

	// Apply manual filters
	query = applyManualFilters(query, accountNumber, search)

	// Debugging query
	query = query.Debug()

	// Initialize paginator
	paginator := pagination.NewPaginator(
		query,
		pagination.WithPage(page),
		pagination.WithPageSize(pageSize),
		pagination.WithSort(sort...),
	)

	// Execute pagination and return results
	var transactions []TransactionData
	result, err := paginator.Paginate(&transactions)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// addDateRangeFilter adds a date range filter to the FilterManager.
func addDateRangeFilter(filterManager *pagination.FilterManager, dateStart, dateEnd string) {
	if dateStart != "" && dateEnd != "" {
		filterManager.AddAndFilter(pagination.DateRangeFilter{
			Field:     "trx_date",
			StartDate: dateStart,
			EndDate:   dateEnd,
		})
	}
}

// addTransactionAmountFilter adds a transaction amount filter to the FilterManager.
func addTransactionAmountFilter(filterManager *pagination.FilterManager, trxAmount float64) {
	if trxAmount > 0 {
		filterManager.AddAndFilter(pagination.ComparisonFilter{
			Field:    "trx_amount",
			Operator: ">=",
			Value:    trxAmount,
		})
	}
}

// applyManualFilters applies manual filters for account number and CIF search.
func applyManualFilters(query *gorm.DB, accountNumber, search string) *gorm.DB {
	if accountNumber != "" {
		query = query.Where("account_number = ?", accountNumber)
	}
	query = query.Where("(trx_type = ? OR trx_type = ?)", "pengeluaran", "pemasukan")
	if search != "" {
		query = query.Where("cif LIKE ?", "%"+search+"%")
	}
	return query
}
