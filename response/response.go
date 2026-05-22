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

func SuccessWithMetaData(ctx *gin.Context, statusCode int, message string, data any, meta any) {
	res := dto.ResponseSuccess{
		Status:  "success",
		Message: message,
	}

	ctx.JSON(statusCode, struct {
		dto.ResponseSuccess
		Data any `json:"data"`
		Meta any `json:"meta"`
	}{
		ResponseSuccess: res,
		Data:            data,
		Meta:            meta,
	})
}

func Error(ctx *gin.Context, statusCodeHTTP int, message string) {
	ctx.JSON(statusCodeHTTP, dto.ResponseError{
		Status:  "error",
		Message: message,
		Error:   http.StatusText(statusCodeHTTP),
	})
}
