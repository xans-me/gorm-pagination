package pagination_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xans-me/gorm-pagination/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type SortTestData struct {
	ID            int
	AccountNumber string
	TrxDate       string
	TrxAmount     float64
}

// Setup test database with test data for sorting tests
func setupSortTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&SortTestData{})

	db.Create(&SortTestData{ID: 1, AccountNumber: "123", TrxDate: "2024-01-01", TrxAmount: 100})
	db.Create(&SortTestData{ID: 2, AccountNumber: "456", TrxDate: "2024-01-02", TrxAmount: 200})
	db.Create(&SortTestData{ID: 3, AccountNumber: "789", TrxDate: "2024-01-03", TrxAmount: 300})

	return db
}

func TestPaginator_Sorting(t *testing.T) {
	db := setupSortTestDB()

	// Test sorting by trx_amount in ascending order
	paginator := pagination.NewPaginator(
		db.Model(&SortTestData{}),
		pagination.WithPage(1),
		pagination.WithPageSize(10),
		pagination.WithSort("trx_amount asc"),
	)

	var results []SortTestData
	res, err := paginator.Paginate(&results)

	assert.Nil(t, err)
	assert.Equal(t, 100.0, results[0].TrxAmount)
	assert.Equal(t, 300.0, results[2].TrxAmount)
	assert.Equal(t, int64(3), res.TotalData)
}
