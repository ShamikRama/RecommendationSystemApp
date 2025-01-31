package erorrs

import "errors"

var (
	ErrRowsNull = errors.New("no rows")
)

var (
	ErrCartNotFound = errors.New("cart not found for user")
	ErrEmptyInput   = errors.New("empty input")
)
