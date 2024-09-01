package pagination

import (
	"gorm.io/gorm"
)

// Pagination represents the general structure for pagination
type Pagination struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	SortField  string `json:"sort_field"`
	SortOrder  string `json:"sort_order"`
	Prefetch   bool   `json:"prefetch"`
	GroupBy    string `json:"group_by"`
}

func NewPagination(page, pageSize int, sortField, sortOrder, groupBy string, prefetch bool) Pagination {
	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		SortField: sortField,
		SortOrder: sortOrder,
		Prefetch:  prefetch,
		GroupBy:   groupBy,
	}
}

// Paginator interface
type Paginator interface {
	Paginate(db *gorm.DB, prefetchSize int) *gorm.DB
	ApplySorting(db *gorm.DB) *gorm.DB
	CalculateTotalPages(totalRows int64)
	PrefetchNextPage(db *gorm.DB)
	GetPrefetchedData() ([]interface{}, error)
	GroupData(db *gorm.DB, groupByField string, filters map[string]interface{}, customSelect, customGroupBy, customOrder string) ([]map[string]interface{}, error)
}

// fetchNextPage as a method on Pagination
func (p *Pagination) fetchNextPage(query *gorm.DB) ([]interface{}, error) {
	var nextData []map[string]interface{}
	if err := query.Find(&nextData).Error; err != nil {
		return nil, err
	}

	interfaceData := make([]interface{}, len(nextData))
	for i, v := range nextData {
		interfaceData[i] = v
	}

	return interfaceData, nil
}

func (p *Pagination) CalculateTotalPages(totalRows int64) {
	p.TotalRows = totalRows
	if totalRows == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = int((totalRows + int64(p.PageSize) - 1) / int64(p.PageSize))
	}
}
