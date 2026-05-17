package appError

import "errors"

var (
	EmailAlreadyExists = errors.New("email already exists")
	InvalidEmailFormat = errors.New("invalid email format")
)
