package pagination

import "errors"

var (
	ErrInvalidPageSize = errors.New("invalid page size, must be greater than 0")
	ErrInvalidPage     = errors.New("invalid page number, must be greater than 0")
)
