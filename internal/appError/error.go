package appError

import "errors"

var (
	EmailAlreadyExists     = errors.New("email already exists")
	InvalidEmailFormat     = errors.New("invalid email format")
	EmailOrPassWrong       = errors.New("wrong email or password")
	UserNotFound           = errors.New("profile not found")
	InvalidPageNumber      = errors.New("invalid page number, must be positive number")
	EmptyPin               = errors.New("pin is empty, create PIN please")
	WrongPin               = errors.New("wrong pin")
	WrongPassword          = errors.New("wrong password")
	PhoneAlreadyExists     = errors.New("phone number already used")
	FileTooLarge           = errors.New("file size exceeds 2MB limit")
	FileTypeNotAllowed     = errors.New("only jpg, png, and webp are allowed")
	InsufficientBalance    = errors.New("insufficient balance")
	ReceiverNotFound       = errors.New("receiver not found")
	SelfTransferNotAllowed = errors.New("cannot transfer to your own account")
	InvalidPeriod          = errors.New("invalid period")
	TokenDoesntExpired     = errors.New("token does not have expiration claim")
	InvalidateSession      = errors.New("failed to invalidate session")
)
