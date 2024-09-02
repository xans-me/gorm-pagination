package pagination

import "gorm.io/gorm"

// Ordering defines an interface for applying orderings.
type Ordering interface {
	Apply(db *gorm.DB) *gorm.DB
}

// OrderBy applies a simple ordering.
type OrderBy struct {
	Field     string
	Direction string // "asc" or "desc"
}

func (o OrderBy) Apply(db *gorm.DB) *gorm.DB {
	return db.Order(o.Field + " " + o.Direction)
}
