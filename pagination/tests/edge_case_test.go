package pagination_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xans-me/gorm-pagination/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type EdgeCaseTestData struct {
	ID            int
	AccountNumber string
	TrxDate       string
	TrxAmount     float64
	TrxType       string
}

// Setup test database with test data for edge cases
func setupEdgeCaseTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&EdgeCaseTestData{})

	return db
}

func TestPaginator_EmptyResults(t *testing.T) {
	db := setupEdgeCaseTestDB()

	paginator := pagination.NewPaginator(
		db.Model(&EdgeCaseTestData{}),
		pagination.WithPage(1),
		pagination.WithPageSize(10),
	)

	var results []EdgeCaseTestData
	res, err := paginator.Paginate(&results)

	assert.Nil(t, err)
	assert.Len(t, results, 0)
	assert.Equal(t, int64(0), res.TotalData)
	assert.Equal(t, 1, res.TotalPages) // Total pages should be 1 even if no data
}
