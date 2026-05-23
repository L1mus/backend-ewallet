package controller

import (
	"strconv"

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
