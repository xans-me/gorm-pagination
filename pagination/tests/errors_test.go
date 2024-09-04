package pagination_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xans-me/gorm-pagination/pagination"
	"testing"
)

func TestErrors(t *testing.T) {
	assert.EqualError(t, pagination.ErrInvalidPageSize, "invalid page size, must be greater than 0")
	assert.EqualError(t, pagination.ErrInvalidPage, "invalid page number, must be greater than 0")
}
