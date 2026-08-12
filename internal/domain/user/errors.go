package user

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("email already registered")
	ErrInvalidCredential = errors.New("email or password wrong")
)
