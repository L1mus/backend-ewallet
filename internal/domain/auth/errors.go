package auth

import "errors"

var (
	ErrInvalidVerificationCode = errors.New("verification code must contain exactly 6 digits")
)
