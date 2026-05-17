package dto

import "time"

type ResponseSuccess struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"errors,omitempty"`
}

type RegisterResponse struct {
	Id        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created-at"`
}
