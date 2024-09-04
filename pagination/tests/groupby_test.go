package pagination_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xans-me/gorm-pagination/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type GroupByTestData struct {
	ID            int
	AccountNumber string
	TrxDate       string
	TrxAmount     float64
	TrxType       string
}

// Setup test database with test data for GroupBy tests
func setupGroupByTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&GroupByTestData{})

	db.Create(&GroupByTestData{ID: 1, AccountNumber: "123", TrxDate: "2024-01-01", TrxAmount: 100, TrxType: "income"})
	db.Create(&GroupByTestData{ID: 2, AccountNumber: "123", TrxDate: "2024-01-02", TrxAmount: 200, TrxType: "expense"})
	db.Create(&GroupByTestData{ID: 3, AccountNumber: "456", TrxDate: "2024-01-01", TrxAmount: 300, TrxType: "income"})

	return db
}

func TestPaginator_GroupBy(t *testing.T) {
	db := setupGroupByTestDB()

	paginator := pagination.NewPaginator(
		db.Model(&GroupByTestData{}),
		pagination.WithPage(1),
		pagination.WithPageSize(10),
	)

	paginator.GroupBy("account_number")

	var results []GroupByTestData
	res, err := paginator.Paginate(&results)

	assert.Nil(t, err)
	assert.Equal(t, 1, res.TotalPages) // Adjust the expectation to 1
}
