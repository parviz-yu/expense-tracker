package errs

import "errors"

var (
	ErrCategoryNotExists = errors.New("category not exists")
	ErrNegativeAmount    = errors.New("negative amount")
)
