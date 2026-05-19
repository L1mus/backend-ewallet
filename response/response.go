package response

import (
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, statusCodeHTTP int, message string, data any) {
	res := dto.ResponseSuccess{
		Status:  "success",
		Message: message,
	}

	if data != nil {
		ctx.JSON(statusCodeHTTP, struct {
			dto.ResponseSuccess
			Data any `json:"data"`
		}{
			ResponseSuccess: res,
			Data:            data,
		})
		return
	}
	ctx.JSON(statusCodeHTTP, res)
}

func Error(ctx *gin.Context, statusCodeHTTP int, message string) {
	ctx.JSON(statusCodeHTTP, dto.ResponseError{
		Status:  "error",
		Message: message,
		Error:   http.StatusText(statusCodeHTTP),
	})
}
