package appError

import "errors"

var (
	EmailAlreadyExists = errors.New("email already exists")
	InvalidEmailFormat = errors.New("invalid email format")
	EmailOrPassWrong   = errors.New("wrong email or password")
	UserNotFound       = errors.New("profile not found")
)
