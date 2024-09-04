package pagination

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// Result contains the paginated result.
type Result struct {
	Data       interface{}            `json:"data"`
	TotalData  int64                  `json:"totalData"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"pageSize"`
	TotalPages int                    `json:"totalPages"`
	Summary    map[string]interface{} `json:"summary,omitempty"`
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
func (p *Paginator) Paginate(result interface{}) (*Result, error) {
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

	return &Result{
		Data:       result,
		TotalData:  p.Total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
		Summary:    summary,
	}, nil
}

// Summary calculates the summary fields dynamically.
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
			// Generic distribution counting based on field value
			var distribution []map[string]interface{}
			p.DB.Model(model).Select(fieldName + ", COUNT(*) as count").Group(fieldName).Order(fieldName).Scan(&distribution)
			summary[fieldName+"_distribution"] = distribution

		case "value_count":
			// Dynamic counting of specific field values
			if len(parts) > 2 {
				// If specific values are provided, split them and count each
				values := strings.Split(parts[2], "|") // Expecting values in format field:aggregationType:value1|value2|...
				for _, value := range values {
					var countResult int64
					p.DB.Model(model).Where(fieldName+" = ?", value).Count(&countResult)
					summary[fieldName+"_"+value+"_count"] = countResult
				}
			} else {
				// If no specific value is provided, count non-NULL values (similar to "count")
				var countResult int64
				p.DB.Model(model).Where(fieldName + " IS NOT NULL").Count(&countResult)
				summary[fieldName+"_count"] = countResult
			}
		}
	}

	return summary
}
