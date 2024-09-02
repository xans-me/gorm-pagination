package pagination

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Paginator handles the pagination logic.
type Paginator struct {
	DB            *gorm.DB
	Page          int
	PageSize      int
	Total         int64
	Sort          []string
	Filters       []Filter
	Groups        []string
	SummaryFields []string
	Orderings     []Ordering
}

// PaginationResult contains the paginated result.
type PaginationResult struct {
	Data       interface{}
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
	Summary    map[string]interface{}
}

// NewPaginator initializes a new Paginator instance with required parameters.
func NewPaginator(db *gorm.DB, options ...PaginatorOption) *Paginator {
	p := &Paginator{
		DB:       db,
		Page:     1,
		PageSize: 10,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// Paginate executes the pagination and returns the result.
func (p *Paginator) Paginate(result interface{}) (*PaginationResult, error) {
	if p.PageSize <= 0 {
		return nil, ErrInvalidPageSize
	}

	if p.Page <= 0 {
		return nil, ErrInvalidPage
	}

	offset := (p.Page - 1) * p.PageSize
	query := p.DB.Offset(offset).Limit(p.PageSize)

	// Apply filters
	for _, filter := range p.Filters {
		query = filter.Apply(query)
	}

	// Apply groupings
	if len(p.Groups) > 0 {
		query = query.Group(clause.GroupBy{
			Columns: make([]clause.Column, len(p.Groups)),
		}.Name(p.Groups...))
	}

	// Apply orderings
	for _, order := range p.Orderings {
		query = order.Apply(query)
	}

	// Apply sorting
	for _, sort := range p.Sort {
		query = query.Order(sort)
	}

	// Fetch paginated results
	if err := query.Find(result).Error; err != nil {
		return nil, err
	}

	// Fetch total count
	if err := p.DB.Model(result).Count(&p.Total).Error; err != nil {
		return nil, err
	}

	// Calculate summary if requested
	summary := p.Summary(result)

	return &PaginationResult{
		Data:       result,
		Total:      p.Total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: int((p.Total + int64(p.PageSize) - 1) / int64(p.PageSize)),
		Summary:    summary,
	}, nil
}

// Summary calculates the summary fields if requested.
func (p *Paginator) Summary(model interface{}) map[string]interface{} {
	if len(p.SummaryFields) == 0 {
		return nil
	}

	summary := make(map[string]interface{})
	for _, field := range p.SummaryFields {
		var result float64
		p.DB.Model(model).Select("SUM(" + field + ")").Scan(&result)
		summary[field] = result
	}
	return summary
}
