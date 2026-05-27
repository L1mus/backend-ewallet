package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/response"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/L1mus/backend-ewallet/pkg"
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

// GetUserDashboard
//
// @Summary      Get user dashboard
// @Description  Get authenticated user's wallet balance, total income, and total expenses
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.ResponseSuccess{data=dto.GetUserDashboardDTO}
// @Failure      500  {object}  dto.ResponseError
// @Router       /users/dashboard [get]
func (c *UserController) GetUserDashboard(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	res, err := c.userService.GetUserDashboard(ctx.Request.Context(), claims.Id)
	if err != nil {
		response.Error(ctx, 500, "internal server error")
	}
	response.Success(ctx, 200, "Get data success", res)
}

// FindReceiver
//
// @Summary      Find transfer receiver
// @Description  Search for other users to be used as a transfer recipient, with pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        search  query     string  false  "Search by name or phone number"  example("John Doe")
// @Param        page    query     string  false  "Page number"                      example("1")
// @Success      200     {object}  dto.FindReceiverResponse
// @Failure      400     {object}  dto.ResponseError
// @Failure      500     {object}  dto.ResponseError
// @Router       /users/transfer [get]
func (c *UserController) FindReceiver(ctx *gin.Context) {
	// bentuk url users?search=&page
	// binding request query
	var req dto.PageQuery
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

// GetTransactionReport
//
// @Summary      Get transaction report
// @Description  Get aggregated income and expense report for the authenticated user, grouped by the selected time period
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        period  query     string  true  "Time period to group by"  Enums(week, month, year)  example("month")
// @Success      200     {object}  dto.GetTransactionReportResponse
// @Failure      400     {object}  dto.ResponseError
// @Failure      500     {object}  dto.ResponseError
// @Router       /users/report [get]
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

// GetTransactionHistory
//
// @Summary      Get transaction history
// @Description  Get a paginated list of the authenticated user's transaction history, optionally filtered by a search keyword
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        search  query     string  false  "Search by receiver name, payment method, or description"  example("John Doe")
// @Param        page    query     string  false  "Page number"                                               example("1")
// @Success      200     {object}  dto.ResponseSuccess{data=[]dto.GetTransactionHistoryDTO}
// @Failure      400     {object}  dto.ResponseError
// @Failure      500     {object}  dto.ResponseError
// @Router       /users/transactions [get]
func (c *UserController) GetTransactionHistory(ctx *gin.Context) {

	var req dto.PageQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "bad request")
	}

	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	res, metaData, err := c.userService.GetTransactionHistory(ctx.Request.Context(), claims.Id, req)
	log.Println("res :\n", res, "meta :\n", metaData, "error :", err)
	if err != nil {
		response.Error(ctx, 500, "internal server")
		return
	}

	response.SuccessWithMetaData(ctx, 200, "Get data success", res, metaData)
}

// UpdateProfile
//
// @Summary      Update User Profile
// @Description  Update full name, phone number, and avatar image simultaneously using multipart/form-data
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        full_name  formData  string  false  "Full Name"
// @Param        phone      formData  string  false  "Phone Number"
// @Param        picture    formData  file    false  "Avatar profile image (jpg, png, max 2MB)"
// @Success      200        {object}  dto.ResponseSuccess
// @Failure      400        {object}  dto.ResponseError
// @Failure      422        {object}  dto.ResponseError
// @Failure      500        {object}  dto.ResponseError
// @Router       /users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var req dto.EditProfileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fileHeader, err := ctx.FormFile("picture")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		response.Error(ctx, http.StatusBadRequest, "Invalid Format")
		return
	}

	err = c.userService.EditProfile(ctx.Request.Context(), claims.Id, req, fileHeader)
	if err != nil {
		if errors.Is(err, appError.PhoneAlreadyExists) {
			response.Error(ctx, 400, err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Profile successfully updated", nil)
}

// EditPin
//
// @Summary      Edit user PIN
// @Description  Update the authenticated user's transaction PIN. Requires the current PIN to be provided
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      dto.EditPinRequest  true  "Edit PIN Payload"
// @Success      200   {object}  dto.ResponseSuccess
// @Failure      400   {object}  dto.ResponseError
// @Failure      500   {object}  dto.ResponseError
// @Router       /users/pin [patch]
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

// EditPassword
//
// @Summary      Edit user password
// @Description  Update the authenticated user's login password. Requires the current password to be provided
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      dto.EditPasswordRequest  true  "Edit Password Payload"
// @Success      200   {object}  dto.ResponseSuccess
// @Failure      400   {object}  dto.ResponseError
// @Failure      500   {object}  dto.ResponseError
// @Router       /users/password [patch]
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
