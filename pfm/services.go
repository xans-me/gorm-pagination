package pfm

import (
	"net/http"
	"strconv"
	"test-pagination-pg-go/pagination"
)

func GetPaginatedTransactions(r *http.Request) (interface{}, error) {
	db := GetDB()

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
	status := r.URL.Query()["status"]

	// Initialize filters
	var filters []pagination.Filter
	if dateStart != "" && dateEnd != "" {
		filters = append(filters, pagination.DateRangeFilter{Field: "date", StartDate: dateStart, EndDate: dateEnd})
	}
	if len(status) > 0 {
		filters = append(filters, pagination.StatusFilter{Field: "status", Statuses: status})
	}

	// Initialize paginator
	paginator := pagination.NewPaginator(
		db,
		pagination.WithPage(page),
		pagination.WithPageSize(pageSize),
		pagination.WithSort(sort...),
		pagination.WithFilters(filters...),
		pagination.WithSummaryFields("trx_amount"),
	)

	var transactions []BrimoPFM
	result, err := paginator.Paginate(&transactions)
	if err != nil {
		return nil, err
	}

	return result, nil
}
