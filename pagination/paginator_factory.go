package pagination

// NewPaginator creates a Paginator instance based on the database type
func NewPaginator(dbType string, page, pageSize int, sortField, sortOrder, groupBy string, prefetch bool) Paginator {
	pagination := NewPagination(page, pageSize, sortField, sortOrder, groupBy, prefetch)

	switch dbType {
	case "postgres":
		return &PostgresPaginator{Pagination: pagination}
	case "mysql":
		return &MySQLPaginator{Pagination: pagination}
	// Add support for other databases here
	default:
		return &PostgresPaginator{Pagination: pagination} // Default to PostgreSQL
	}
}
