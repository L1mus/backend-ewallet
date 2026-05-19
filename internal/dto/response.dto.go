package dto

import (
	"time"
)

type ResponseSuccess struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Welcome, John doe"`
}

type ResponseError struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed get data"`
	Error   string `json:"errors"`
}

type RegisterResponse struct {
	Id        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created-at"`
}

type LoginResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type GetUserProfileResponse struct {
	Id                int     `json:"id"`
	FullName          string  `json:"full_name"`
	Email             string  `json:"email"`
	Phone             *string `json:"phone"`
	ProfilePictureURL *string `json:"profile_picture_url"`
}
