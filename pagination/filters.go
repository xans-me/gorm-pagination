package pagination

import (
	"gorm.io/gorm"
	"strings"
)

// Filter defines an interface for applying filters.
type Filter interface {
	Apply(db *gorm.DB) *gorm.DB
}

// FilterManager manages multiple filters and applies them.
type FilterManager struct {
	AndFilters []Filter
	OrFilters  []Filter
}

// AddAndFilter adds a filter that will be combined with AND logic.
func (fm *FilterManager) AddAndFilter(filter Filter) {
	fm.AndFilters = append(fm.AndFilters, filter)
}

// AddOrFilter adds a filter that will be combined with OR logic.
func (fm *FilterManager) AddOrFilter(filter Filter) {
	fm.OrFilters = append(fm.OrFilters, filter)
}

// Apply applies all the filters in the FilterManager to the query.
func (fm *FilterManager) Apply(db *gorm.DB) *gorm.DB {
	// Apply AND filters
	for _, filter := range fm.AndFilters {
		db = filter.Apply(db)
	}

	// Build OR filters using raw SQL
	if len(fm.OrFilters) > 0 {
		var orConditions []string
		var orValues []interface{}

		// Build the OR conditions as raw SQL
		for _, filter := range fm.OrFilters {
			query, args := filterToSQL(filter)
			orConditions = append(orConditions, query)
			orValues = append(orValues, args...)
		}

		// Combine OR conditions
		if len(orConditions) > 0 {
			orClause := "(" + strings.Join(orConditions, " OR ") + ")"
			db = db.Where(orClause, orValues...)
		}
	}

	return db
}

// Helper function to convert Filter to raw SQL
func filterToSQL(filter Filter) (string, []interface{}) {
	switch f := filter.(type) {
	case ComparisonFilter:
		return f.Field + " " + f.Operator + " ?", []interface{}{f.Value}
	default:
		return "", nil
	}
}

// DateRangeFilter applies a date range filter.
type DateRangeFilter struct {
	Field     string
	StartDate string
	EndDate   string
}

func (f DateRangeFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where(f.Field+" BETWEEN ? AND ?", f.StartDate, f.EndDate)
}

// ComparisonFilter allows filtering with different comparison operators.
type ComparisonFilter struct {
	Field    string
	Operator string // Examples: "=", ">", "<", ">=", "<=", "!="
	Value    interface{}
}

func (f ComparisonFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where(f.Field+" "+f.Operator+" ?", f.Value)
}

// StatusFilter applies a status filter (used as an example of IN clause).
type StatusFilter struct {
	Field    string
	Statuses []string
}

func (f StatusFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where(f.Field+" IN ?", f.Statuses)
}

// SearchFilter applies a search filter using LIKE (for string searches).
type SearchFilter struct {
	Field string
	Value string
}

func (f SearchFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where(f.Field+" LIKE ?", "%"+f.Value+"%")
}
