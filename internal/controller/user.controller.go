package controller

import (
	"strconv"

	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/L1mus/backend-ewallet/response"
	"github.com/gin-gonic/gin"
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
	res, err := c.userService.GetUserDashboad(ctx, claims.Id)
	if err != nil {
		response.Error(ctx, 500, "internal server error")
	}
	response.Success(ctx, 200, "Get data success", res)
}

func (c *UserController) FindReceiver(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	search := ctx.Query("search")
	strPage := ctx.DefaultQuery("page", "1")
	strLimit := ctx.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(strPage)
	limit, _ := strconv.Atoi(strLimit)
	res, err := c.userService.FindReceiver(ctx, claims.Id, search, limit, page)
	if err != nil {
		response.Error(ctx, 500, "internal server error")
	}
	response.Success(ctx, 200, "Get data success", res)
}
