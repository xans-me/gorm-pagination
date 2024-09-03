package pfm

import (
	"net/http"
	"strconv"
	"strings"
	"test-pagination-pg-go/pagination"
)

func GetPaginatedTransactions(r *http.Request) (interface{}, error) {
	db := GetDB()
	db.Debug()

	// Get query parameters
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
	dateColumn := r.URL.Query().Get("dateColumn")
	if dateColumn == "" {
		dateColumn = "trx_date" // Default date column
	}
	status := r.URL.Query()["status"]
	trxType := r.URL.Query().Get("trx_type") // Get the transaction type filter

	// Summary fields can be specified as comma-separated values
	summaryFields := r.URL.Query().Get("summaryFields")
	var summaryFieldList []string
	if summaryFields != "" {
		summaryFieldList = strings.Split(summaryFields, ",")
	} else {
		summaryFieldList = []string{} // No default summary field; client must specify
	}

	// Start building the query
	query := db.Model(&BrimoPFM{})

	// Apply date range filter
	if dateStart != "" && dateEnd != "" {
		query = query.Where(dateColumn+" BETWEEN ? AND ?", dateStart, dateEnd)
	}

	// Apply status filter
	if len(status) > 0 {
		query = query.Where("status = ?", status)
	}

	// Apply transaction type filter
	if trxType != "" {
		query = query.Where("trx_type = ?", trxType)
	}

	// Initialize paginator with the built query
	paginator := pagination.NewPaginator(
		query, // Pass the query object instead of db
		pagination.WithPage(page),
		pagination.WithPageSize(pageSize),
		pagination.WithSort(sort...),
		pagination.WithSummaryFields(summaryFieldList...),
	)

	var transactions []BrimoPFM
	result, err := paginator.Paginate(&transactions)
	if err != nil {
		return nil, err
	}

	return result, nil
}
