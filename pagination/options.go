package pagination

// PaginatorOption is a function that configures a Paginator.
type PaginatorOption func(*Paginator)

// WithPage sets the current page for the paginator.
func WithPage(page int) PaginatorOption {
	return func(p *Paginator) {
		if page > 0 {
			p.Page = page
		}
	}
}

// WithPageSize sets the page size for the paginator.
func WithPageSize(pageSize int) PaginatorOption {
	return func(p *Paginator) {
		if pageSize > 0 {
			p.PageSize = pageSize
		}
	}
}

// WithSort sets the sorting options for the paginator.
func WithSort(sort ...string) PaginatorOption {
	return func(p *Paginator) {
		p.Sort = sort
	}
}

// WithSummaryFields sets the fields for which summaries should be calculated.
func WithSummaryFields(fields ...string) PaginatorOption {
	return func(p *Paginator) {
		p.SummaryFields = fields
	}
}
