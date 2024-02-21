package errs

import "errors"

var (
	ErrCategoryNotExists   = errors.New("category not exists")
	ErrNegativeAmount      = errors.New("negative amount")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrInvalidQueryParams  = errors.New("invalid query params")
	ErrInternalServerError = errors.New("internal server error")
)
