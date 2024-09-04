package pagination_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xans-me/gorm-pagination/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type FilterTestData struct {
	ID    int
	Name  string
	Age   int
	Email string
}

// Setup test database with test data
func setupFilterTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&FilterTestData{})

	db.Create(&FilterTestData{ID: 1, Name: "John", Age: 30, Email: "john@example.com"})
	db.Create(&FilterTestData{ID: 2, Name: "Jane", Age: 25, Email: "jane@example.com"})
	db.Create(&FilterTestData{ID: 3, Name: "Doe", Age: 35, Email: "doe@example.com"})

	return db
}

func TestComparisonFilter(t *testing.T) {
	db := setupFilterTestDB()

	filter := pagination.ComparisonFilter{
		Field:    "age",
		Operator: ">=",
		Value:    30,
	}

	query := filter.Apply(db.Model(&FilterTestData{}))

	var results []FilterTestData
	err := query.Find(&results).Error

	assert.Nil(t, err)
	assert.Len(t, results, 2)
}

func TestSearchFilter(t *testing.T) {
	db := setupFilterTestDB()

	filter := pagination.SearchFilter{
		Field: "email",
		Value: "john",
	}

	query := filter.Apply(db.Model(&FilterTestData{}))

	var results []FilterTestData
	err := query.Find(&results).Error

	assert.Nil(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "john@example.com", results[0].Email)
}
