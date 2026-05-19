package controller

import (
	"fmt"
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/L1mus/backend-ewallet/pkg"
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
// @Accept       JSON
// @Produce      JSON
// @Param		 claims.id	claims.id	integer true
// @Success      200  {object}  dto.ResponseSuccess
// @Failure      500  {object}  dto.ResponseError
// @Router       /users/profile [post]

func (c *UserController) GetUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	fmt.Println(token)
	fmt.Println("Claims id", claims.Id)
	response, err := c.userService.GetUserProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  "error",
			Message: err.Error(),
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  "success",
		Message: "Get data success",
		Data:    response,
	})
}
