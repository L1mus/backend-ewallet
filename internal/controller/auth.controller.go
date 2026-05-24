package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/response"
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

// Register
//
// @Summary      Register account
// @Description  create user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param		 body	body	dto.RegisterRequest true "Register Payload"
// @Success      201  {object}  dto.RegisterResponse
// @Failure      400  {object}  dto.ResponseError
// @Failure      500  {object}  dto.ResponseError
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterRequest

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {

		if strings.Contains(err.Error(), "FullName") && strings.Contains(err.Error(), "required") {
			response.Error(ctx, 400, "FullName is required")
			return
		}
		if strings.Contains(err.Error(), "Email") && strings.Contains(err.Error(), "required") {
			response.Error(ctx, 400, "Email is required")
			return
		}
		if strings.Contains(err.Error(), "RegisterRequest.Password") && strings.Contains(err.Error(), "required") {
			response.Error(ctx, 400, "Password is required")
			return
		}
		if strings.Contains(err.Error(), "RegisterRequest.ConfirmPassword") && strings.Contains(err.Error(), "required") {
			response.Error(ctx, 400, "Confirm Password is required")
			return
		}
		if strings.Contains(err.Error(), "ConfirmPassword") {
			response.Error(ctx, 400, "Password confirmation does not match")
			return
		}
		if strings.Contains(err.Error(), "email") && strings.Contains(err.Error(), "validation") {
			response.Error(ctx, 400, "invalid email format")
			return
		}
		response.Error(ctx, 500, err.Error())
		return
	}

	res, err := c.authService.Register(ctx.Request.Context(), body)
	if err != nil {
		if errors.Is(err, appError.EmailAlreadyExists) || errors.Is(err, appError.InvalidEmailFormat) {
			response.Error(ctx, 400, err.Error())
		} else {
			response.Error(ctx, 500, err.Error())
		}
		return
	}
	response.Success(ctx, 201, fmt.Sprintf("Register Complete, Welcome %s", res.FullName), res)
}

// Login
//
// @Summary      Login account
// @Description  login into user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param		 body	body	dto.LoginRequest true "Login Payload"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  dto.ResponseError
// @Failure      500  {object}  dto.ResponseError
// @Router       /auth [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var body dto.LoginRequest

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}
	res, err := c.authService.Login(ctx.Request.Context(), body)
	if err != nil {
		response.Error(ctx, 500, "internal server error")
		return
	}
	response.Success(ctx, 200, fmt.Sprintf("Login Complete, Welcome %s", res.FullName), res)
}
