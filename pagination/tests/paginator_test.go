package pagination_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xans-me/gorm-pagination/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type TestData struct {
	ID            int
	AccountNumber string
	TrxDate       string
	TrxAmount     float64
	TrxType       string
	CIF           string
}

// Setup test database with test data
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&TestData{})

	db.Create(&TestData{ID: 1, AccountNumber: "123", TrxDate: "2024-01-01", TrxAmount: 100, TrxType: "income", CIF: "ABC123"})
	db.Create(&TestData{ID: 2, AccountNumber: "456", TrxDate: "2024-02-01", TrxAmount: 200, TrxType: "expense", CIF: "DEF456"})
	db.Create(&TestData{ID: 3, AccountNumber: "789", TrxDate: "2024-03-01", TrxAmount: 300, TrxType: "income", CIF: "GHI789"})

	return db
}

func TestPaginator_Paginate(t *testing.T) {
	db := setupTestDB()

	paginator := pagination.NewPaginator(
		db.Model(&TestData{}),
		pagination.WithPage(1),
		pagination.WithPageSize(2),
	)

	var results []TestData
	res, err := paginator.Paginate(&results)

	assert.Nil(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(3), res.TotalData)
	assert.Equal(t, 2, res.TotalPages)
}

func TestPaginator_Filtering(t *testing.T) {
	db := setupTestDB()

	filterManager := pagination.FilterManager{}
	filterManager.AddAndFilter(pagination.ComparisonFilter{
		Field:    "trx_type",
		Operator: "=",
		Value:    "income",
	})

	query := filterManager.Apply(db.Model(&TestData{}))

	paginator := pagination.NewPaginator(
		query,
		pagination.WithPage(1),
		pagination.WithPageSize(10),
	)

	var results []TestData
	res, err := paginator.Paginate(&results)

	assert.Nil(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "income", results[0].TrxType)
	assert.Equal(t, int64(2), res.TotalData)
}

func TestPaginator_Summary(t *testing.T) {
	db := setupTestDB()

	paginator := pagination.NewPaginator(
		db.Model(&TestData{}),
		pagination.WithPage(1),
		pagination.WithPageSize(10),
		pagination.WithSummaryFields("trx_amount:sum", "trx_amount:min", "trx_amount:max"),
	)

	var results []TestData
	res, err := paginator.Paginate(&results)

	assert.Nil(t, err)
	assert.NotNil(t, res.Summary)
	assert.Equal(t, float64(100), res.Summary["trx_amount_min"])
	assert.Equal(t, float64(300), res.Summary["trx_amount_max"])
	assert.Equal(t, float64(600), res.Summary["trx_amount_sum"])
}
