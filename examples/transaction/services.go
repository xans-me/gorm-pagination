package transaction

import (
	"github.com/xans-me/gorm-pagination/pagination"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	query := db.Model(&Data{})

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

	// Initialize paginator with summary fields for trx_amount, trx_type, and dynamic counting
	paginator := pagination.NewPaginator(
		query,
		pagination.WithPage(page),
		pagination.WithPageSize(pageSize),
		pagination.WithSort(sort...),
		// Adding various summary fields dynamically
		pagination.WithSummaryFields(
			"trx_amount:sum",
			"trx_amount:min",                      // Min of trx_amount
			"trx_amount:max",                      // Max of trx_amount// Sum of trx_amount
			"account_number:distribution",         // Distribution of trx_type
			"trx_type:value_count:income|expense", // Count of 'income' and 'expense'
			"trx_type:value_count"),               // This will count all non-NULL trx_type records
	)

	// Execute pagination and return results
	var transactions []Data
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
	query = query.Where("(trx_type = ? OR trx_type = ?)", "expense", "income")
	if search != "" {
		query = query.Where("cif LIKE ?", "%"+search+"%")
	}
	return query
}
