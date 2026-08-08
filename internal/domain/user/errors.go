package user

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("email already registered")
)
