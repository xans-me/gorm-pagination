package pfm

import (
	"gorm.io/gorm"
)

type BrimoPFMRepository interface {
	FindAll(db *gorm.DB, filters map[string]interface{}) (*gorm.DB, int64, error)
	GroupByPeriod(db *gorm.DB, groupByField string, filters map[string]interface{}, customSelect, customGroupBy, customOrder string) ([]map[string]interface{}, error)
}

type brimoPFMRepository struct{}

func NewBrimoPFMRepository() BrimoPFMRepository {
	return &brimoPFMRepository{}
}

func (r *brimoPFMRepository) FindAll(db *gorm.DB, filters map[string]interface{}) (*gorm.DB, int64, error) {
	query := db.Model(&BrimoPFM{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	var totalRows int64
	query.Count(&totalRows)

	return query, totalRows, nil
}

func (r *brimoPFMRepository) GroupByPeriod(db *gorm.DB, groupByField string, filters map[string]interface{}, customSelect, customGroupBy, customOrder string) ([]map[string]interface{}, error) {
	// Implementasi group data berdasarkan period
	var groupedData []map[string]interface{}
	groupQuery := db.Model(&BrimoPFM{}).
		Select(customSelect).
		Group(customGroupBy).
		Order(customOrder)

	for key, value := range filters {
		groupQuery = groupQuery.Where(key+" = ?", value)
	}

	err := groupQuery.Scan(&groupedData).Error
	return groupedData, err
}
