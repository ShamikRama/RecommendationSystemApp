package erorrs

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRows           = sql.ErrNoRows
	ErrUniqueConstraint = errors.New("unique constraint violation")
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmptyInput         = errors.New("input len is null")
)
