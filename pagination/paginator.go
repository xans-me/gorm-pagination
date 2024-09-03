package pagination

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
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
		groupByClause := clause.GroupBy{
			Columns: make([]clause.Column, len(p.Groups)),
		}
		for i, group := range p.Groups {
			groupByClause.Columns[i] = clause.Column{Name: group}
		}
		query = query.Clauses(groupByClause)
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

	// Calculate TotalPages safely
	totalPages := int(p.Total / int64(p.PageSize))
	if p.Total%int64(p.PageSize) != 0 {
		totalPages++
	}

	// Calculate summary if requested
	summary := p.Summary(result)

	return &PaginationResult{
		Data:       result,
		Total:      p.Total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
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
		// Expecting field to be in format "field:aggregationType"
		parts := strings.Split(field, ":")
		fieldName := parts[0]
		aggregationType := "sum" // Default to sum if not specified
		if len(parts) > 1 {
			aggregationType = parts[1]
		}

		switch aggregationType {
		case "sum":
			var sumResult float64
			p.DB.Model(model).Select("SUM(" + fieldName + ")").Scan(&sumResult)
			summary[fieldName+"_sum"] = sumResult

		case "min":
			var minResult float64
			p.DB.Model(model).Select("MIN(" + fieldName + ")").Scan(&minResult)
			summary[fieldName+"_min"] = minResult

		case "max":
			var maxResult float64
			p.DB.Model(model).Select("MAX(" + fieldName + ")").Scan(&maxResult)
			summary[fieldName+"_max"] = maxResult

		case "distribution":
			var distribution []map[string]interface{}
			p.DB.Model(model).Select(fieldName + ", COUNT(*) as count").Group(fieldName).Order(fieldName).Scan(&distribution)
			summary[fieldName+"_distribution"] = distribution

		case "top":
			// Top N Categories (e.g., top 5 categories)
			topN := 5 // Default to top 5, can be parameterized if needed
			var topCategories []map[string]interface{}
			p.DB.Model(model).Select(fieldName + ", COUNT(*) as count").Group(fieldName).Order("count DESC").Limit(topN).Scan(&topCategories)
			summary[fieldName+"_top_"+strconv.Itoa(topN)+"_categories"] = topCategories
		}
	}

	return summary
}
