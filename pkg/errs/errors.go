package errs

import "errors"

var (
	ErrCategoryNotExists    = errors.New("category not exists")
	ErrNegativeAmount       = errors.New("negative amount")
	ErrInvalidRequestBody   = errors.New("invalid request body")
	ErrInternalServerError  = errors.New("internal server error")
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrInvalidDateFormat    = errors.New("invalid date format")
	ErrInvalidMinMaxAmounts = errors.New("invalid min max amounts")
)
