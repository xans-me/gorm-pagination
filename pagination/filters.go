package pagination

import "gorm.io/gorm"

// Filter defines an interface for applying filters.
type Filter interface {
	Apply(db *gorm.DB) *gorm.DB
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

// StatusFilter applies a status filter.
type StatusFilter struct {
	Field    string
	Statuses []string
}

func (f StatusFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where(f.Field+" IN ?", f.Statuses)
}

// SearchFilter applies a search filter.
type SearchFilter struct {
	Field string
	Value string
}

func (f SearchFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where(f.Field+" LIKE ?", "%"+f.Value+"%")
}
