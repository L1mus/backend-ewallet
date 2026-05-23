package controller

import (
	"errors"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/L1mus/backend-ewallet/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUserProfile
//
// @Summary      Get user Profile
// @Description  Get user Profile
// @Tags         users
// @Accept       json
// @Produce      json
// @Security	 APIKeyAuth
// @Success      200  {object}  dto.GetUserProfileResponse
// @Failure      500  {object}  dto.ResponseError
// @Router       /users/profile [get]
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	res, err := c.userService.GetUserProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}
	response.Success(ctx, 200, "Get data success", res)
}

func (c *UserController) GetUserDashboard(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	res, err := c.userService.GetUserDashboard(ctx.Request.Context(), claims.Id)
	if err != nil {
		response.Error(ctx, 500, "internal server error")
	}
	response.Success(ctx, 200, "Get data success", res)
}

func (c *UserController) FindReceiver(ctx *gin.Context) {
	// bentuk url users?search=&page
	// binding request query
	var req dto.ReceiverQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "bad request")
	}
	// default page jika request page tidak string kosong
	// get all data, get total pages
	// ambil seluruh data dengan mengirim request sebagai params
	// inisialisasi nextpage dan prevpage
	// format page string dengan Sprintf(/users?search=&page=%d%s)
	// response metadata

	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	//search := ctx.Query("search")
	//strPage := ctx.DefaultQuery("page", "1")
	//strLimit := ctx.DefaultQuery("limit", "10")
	//page, _ := strconv.Atoi(strPage)
	//limit, _ := strconv.Atoi(strLimit)

	res, metaData, err := c.userService.FindReceiver(ctx.Request.Context(), claims.Id, req)
	if err != nil {
		response.Error(ctx, 500, "internal server")
		return
	}

	response.SuccessWithMetaData(ctx, 200, "Get data success", res, metaData)
}

func (c *UserController) GetTransactionReport(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	var body dto.GetTransactionsReportRequest
	if err := ctx.ShouldBindQuery(&body); err != nil {
		response.Error(ctx, 400, "bad request")
		return
	}
	res, err := c.userService.GetTransactionReport(ctx.Request.Context(), claims.Id, body.Period)
	if err != nil {
		response.Error(ctx, 500, "internal server error")
		return
	}
	response.Success(ctx, 200, "Get data success", res)
}

func (c *UserController) GetTransactionHistory(ctx *gin.Context) {

	var req dto.TransactionHistoryQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "bad request")
	}

	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	res, metaData, err := c.userService.GetTransactionHistory(ctx.Request.Context(), claims.Id, req)
	if err != nil {
		response.Error(ctx, 500, "internal server")
		return
	}

	response.SuccessWithMetaData(ctx, 200, "Get data success", res, metaData)
}

func (c *UserController) EditProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.EditProfileRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}

	err := c.userService.EditProfile(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		if errors.Is(err, appError.PhoneAlreadyExists) {
			response.Error(ctx, 400, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}

	response.Success(ctx, 200, "Profile updated successfully", nil)
}

func (c *UserController) EditPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.EditPinRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}

	err := c.userService.EditPin(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		if errors.Is(err, appError.WrongPin) || errors.Is(err, appError.EmptyPin) {
			response.Error(ctx, 400, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}

	response.Success(ctx, 200, "PIN updated successfully", nil)
}

func (c *UserController) EditPassword(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.EditPasswordRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}

	err := c.userService.EditPassword(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		if errors.Is(err, appError.WrongPassword) {
			response.Error(ctx, 400, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}

	response.Success(ctx, 200, "Password updated successfully", nil)
}

func (c *UserController) UploadProfilePicture(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	fileHeader, err := ctx.FormFile("picture")
	if err != nil {
		response.Error(ctx, 400, "picture file is required")
		return
	}

	url, err := c.userService.UploadProfilePicture(ctx.Request.Context(), claims.Id, fileHeader)
	if err != nil {
		if errors.Is(err, appError.FileTooLarge) || errors.Is(err, appError.FileTypeNotAllowed) {
			response.Error(ctx, 422, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}

	response.Success(ctx, 200, "Profile picture updated", gin.H{
		"url": url,
	})
}
