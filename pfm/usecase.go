package pfm

import (
	"test-pagination-pg-go/pagination"

	"gorm.io/gorm"
)

type BrimoPFMService interface {
	GetBrimoPFM(db *gorm.DB, page, pageSize int, sortField, sortOrder, groupBy string, filters map[string]interface{}, prefetch bool) (map[string]interface{}, error)
}

type brimoPFMService struct {
	repo BrimoPFMRepository
}

func NewBrimoPFMService(repo BrimoPFMRepository) BrimoPFMService {
	return &brimoPFMService{repo}
}

func (s *brimoPFMService) GetBrimoPFM(db *gorm.DB, page, pageSize int, sortField, sortOrder, groupBy string, filters map[string]interface{}, prefetch bool) (map[string]interface{}, error) {
	// Initialize Paginator
	paginator := pagination.NewPaginator("postgres", page, pageSize, sortField, sortOrder, groupBy, prefetch)

	var totalAmountOverall float64

	// Calculate total amount of all data matching filters
	db.Model(&BrimoPFM{}).Where(filters).Select("SUM(trx_amount)").Scan(&totalAmountOverall)

	if groupBy != "" {
		groupedData, err := paginator.GroupData(db, groupBy, filters, "", "", "")
		if err != nil {
			return nil, err
		}

		// Return grouped_data response
		return map[string]interface{}{
			"pagination": map[string]interface{}{
				"page":        paginator.(*pagination.PostgresPaginator).Page,
				"page_size":   paginator.(*pagination.PostgresPaginator).PageSize,
				"total_rows":  paginator.(*pagination.PostgresPaginator).TotalRows,
				"total_pages": paginator.(*pagination.PostgresPaginator).TotalPages,
			},
			"grouped_data":         groupedData,
			"total_amount_overall": totalAmountOverall,
		}, nil
	}

	query, totalRows, err := s.repo.FindAll(db, filters)
	if err != nil {
		return nil, err
	}

	// Paginate and fetch data
	query = paginator.Paginate(query, pageSize)

	var data []BrimoPFM
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	// Return data response
	return map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"page":        paginator.(*pagination.PostgresPaginator).Page,
			"page_size":   paginator.(*pagination.PostgresPaginator).PageSize,
			"total_rows":  totalRows,
			"total_pages": paginator.(*pagination.PostgresPaginator).TotalPages,
		},
		"total_amount_overall": totalAmountOverall,
	}, nil
}
