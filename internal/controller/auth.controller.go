package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error : ", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  "Error",
			Message: err.Error(),
			Error:   "Bad Request",
		})
		return
	}

	response, err := c.authService.Register(ctx.Request.Context(), body)
	if err != nil {
		if errors.Is(err, appError.EmailAlreadyExists) || errors.Is(err, appError.InvalidEmailFormat) {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Status:  "Error",
				Message: err.Error(),
				Error:   "Bad Request",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
				Status:  "Error",
				Message: "Internal Server Error",
			})
		}
		return
	}
	ctx.JSON(http.StatusCreated, dto.ResponseSuccess{
		Status:  "Success",
		Message: fmt.Sprintf("Register Complete, Welcome %s", response.FullName),
		Data:    response,
	})
}
