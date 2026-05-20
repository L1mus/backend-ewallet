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
	Message string `json:"message" example:"Failed get data/internal server error"`
	Error   string `json:"errors" example:"internal server error/bad request"`
}

type RegisterResponse struct {
	ResponseSuccess
	Data RegisterDTO
}

type LoginResponse struct {
	ResponseSuccess
	Data LoginDTO
}

type RegisterDTO struct {
	Id        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created-at"`
}

type LoginDTO struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type GetUserProfileDTO struct {
	Id                int     `json:"id" example:"1"`
	FullName          string  `json:"full_name" example:"John Doe"`
	Email             string  `json:"email" example:"example@mail.com"`
	Phone             *string `json:"phone" example:"021234512552"`
	ProfilePictureURL *string `json:"profile_picture_url" example:"https://example.com"`
}
