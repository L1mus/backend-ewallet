package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/response"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/L1mus/backend-ewallet/pkg"
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
// @Failure      400  {object}  dto.ResponseError "required field || invalid format || email already exist"
// @Failure      409  {object}  dto.ResponseError
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
		if errors.Is(err, appError.InvalidEmailFormat) {
			response.Error(ctx, 400, err.Error())
		} else if errors.Is(err, appError.EmailAlreadyExists) {
			response.Error(ctx, 409, err.Error())
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
// @Failure      400  {object}  dto.ResponseError "Email or pass wrong || binding error"
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
		if errors.Is(err, appError.EmailOrPassWrong) {
			response.Error(ctx, 400, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}
	response.Success(ctx, 200, fmt.Sprintf("Login Complete, Welcome %s", res.FullName), res)
}

// Logout
//
// @Summary      Logout account
// @Description  Blacklist current active session token
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.ResponseSuccess
// @Failure      422  {object}  dto.ResponseError
// @Failure      500  {object}  dto.ResponseError
// @Router       /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	/*
		Ambil claims
		lihat kapan token kadaluarsa
		ambil token mentah
		Ambil waktu expired
		untuk melihat Sisa waktu hidup token dijadikan durasi TTL di Redis
		insert token ke redis

	*/
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	tokenStr, _ := ctx.Get("raw_token")
	err := c.authService.Logout(ctx.Request.Context(), claims, tokenStr.(string))
	if err != nil {
		if errors.Is(err, appError.TokenDoesntExpired) {
			response.Error(ctx, 422, err.Error())
			return
		}
		if errors.Is(err, appError.InvalidateSession) {
			response.Error(ctx, 500, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}
	response.Success(ctx, 200, "Logout complete, session end", nil)
}

// ForgotPassword
//
// @Summary      Forgot password
// @Description  Request a password reset token. In production this would be sent via email. Token is valid for 1 hour.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.ForgotPasswordRequest  true  "Forgot Password Payload"
// @Success      200   {object}  dto.ResponseSuccess{data=dto.ForgotPasswordDTO}
// @Failure      400   {object}  dto.ResponseError  "invalid email format"
// @Failure      404   {object}  dto.ResponseError  "user not found"
// @Failure      500   {object}  dto.ResponseError  "internal server error"
// @Router       /auth/forgot-password [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var body dto.ForgotPasswordRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}

	res, err := c.authService.ForgotPassword(ctx.Request.Context(), body)
	if err != nil {
		if errors.Is(err, appError.UserNotFound) {
			response.Error(ctx, 404, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}

	response.Success(ctx, 200, "Reset token generated successfully", res)
}

// ResetPassword
//
// @Summary      Reset password
// @Description  Reset user password using a valid token from forgot-password endpoint. Token can only be used once and expires in 1 hour.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.ResetPasswordRequest  true  "Reset Password Payload"
// @Success      200   {object}  dto.ResponseSuccess
// @Failure      400   {object}  dto.ResponseError  "Validate Binding GIN"
// @Failure      404   {object}  dto.ResponseError  "token not found"
// @Failure      422   {object}  dto.ResponseError  "token used | token expired"
// @Failure      500   {object}  dto.ResponseError  "internal server error"
// @Router       /auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var body dto.ResetPasswordRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}

	err := c.authService.ResetPassword(ctx.Request.Context(), body)
	if err != nil {
		if errors.Is(err, appError.ForgotPasswordTokenNotFound) {
			response.Error(ctx, 404, err.Error())
			return
		}
		if errors.Is(err, appError.ForgotPasswordTokenExpired) ||
			errors.Is(err, appError.ForgotPasswordTokenUsed) {
			response.Error(ctx, 422, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}

	response.Success(ctx, 200, "Password reset successfully", nil)
}
