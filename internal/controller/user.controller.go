package controller

import (
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
