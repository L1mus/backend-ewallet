package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(ctx *gin.Context) {
	// mengambil token dari request payload (header)
	// header Authorization
	// Bearer token => token diawali dengan kata Bearer
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
			Status:  "Unauthorized",
			Message: "Unauthorized Access, Please Login",
			Error:   "unauthorized access, please login",
		})
		return
	}
	splittedBearer := strings.Split(bearerToken, " ")
	if len(splittedBearer) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
			Status:  "Unauthorized",
			Message: "Unauthorized Access, Please Login",
			Error:   "invalid token",
		})
		return
	}
	token := splittedBearer[1]

	// verifikasi token nya
	var claims pkg.Claims
	if err := claims.VerifyJWT(token); err != nil {
		log.Println("Error: ", err.Error())
		if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenExpired) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
				Message: "Unauthorized Access, Please Login",
				Status:  "Unauthorized",
				Error:   err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError{
			Message: "internal server error",
			Status:  "error",
			Error:   "internal server error",
		})
		return
	}

	// menempelkan (attach) claims ke context request
	ctx.Set("claims", claims)
	ctx.Next()
}
