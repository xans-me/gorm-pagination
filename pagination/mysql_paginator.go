package pagination

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

// MySQLPaginator handles pagination for MySQL
type MySQLPaginator struct {
	Pagination
	NextPageData  []interface{}
	NextPageError error
	prefetchLock  sync.Mutex
}

func (p *MySQLPaginator) Paginate(db *gorm.DB, prefetchSize int) *gorm.DB {
	offset := (p.Page - 1) * p.PageSize
	return db.Offset(offset).Limit(p.PageSize + prefetchSize)
}

func (p *MySQLPaginator) ApplySorting(db *gorm.DB) *gorm.DB {
	if p.SortField != "" && (p.SortOrder == "asc" || p.SortOrder == "desc") {
		db = db.Order(fmt.Sprintf("%s %s", p.SortField, p.SortOrder))
	}
	return db
}

func (p *MySQLPaginator) PrefetchNextPage(db *gorm.DB) {
	if !p.Prefetch || p.Page >= p.TotalPages {
		return
	}

	p.prefetchLock.Lock()
	defer p.prefetchLock.Unlock()

	go func() {
		offset := p.Page * p.PageSize
		nextPageQuery := db.Offset(offset).Limit(p.PageSize)
		p.NextPageData, p.NextPageError = p.fetchNextPage(nextPageQuery)
	}()
}

func (p *MySQLPaginator) GetPrefetchedData() ([]interface{}, error) {
	p.prefetchLock.Lock()
	defer p.prefetchLock.Unlock()

	return p.NextPageData, p.NextPageError
}

func (p *MySQLPaginator) GroupData(db *gorm.DB, groupByField string, filters map[string]interface{}, customSelect, customGroupBy, customOrder string) ([]map[string]interface{}, error) {
	// Jika customSelect tidak disediakan, gunakan SELECT default
	if customSelect == "" {
		customSelect = groupByField + ", COUNT(*) AS count"
	}

	// Jika customGroupBy tidak disediakan, gunakan GROUP BY default
	if customGroupBy == "" {
		customGroupBy = groupByField
	}

	// Jika customOrder tidak disediakan, gunakan ORDER BY default
	if customOrder == "" {
		customOrder = groupByField
	}

	// Buat query dengan kustomisasi SELECT, GROUP BY, dan ORDER BY
	groupQuery := db.Model(&map[string]interface{}{}).
		Select(customSelect).
		Group(customGroupBy).
		Order(customOrder)

	// Terapkan filter dinamis
	for key, value := range filters {
		groupQuery = groupQuery.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Jalankan query dan pindahkan hasil ke map
	var groupedData []map[string]interface{}
	if err := groupQuery.Scan(&groupedData).Error; err != nil {
		return nil, err
	}

	return groupedData, nil
}
