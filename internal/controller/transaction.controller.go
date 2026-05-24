package controller

import (
	"errors"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/response"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// CreateTransfer
//
// @Summary      Create money transfer
// @Description  Transfer funds to another user. Requires a PIN for confirmation. The sender's wallet balance decreases and the recipient's wallet balance increases atomically (DB transaction). Transfers to oneself are not permitted.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      dto.CreateTransferRequest  true  "Transfer Payload"
// @Success      201   {object}  dto.ResponseSuccess
// @Failure      400   {object}  dto.ResponseError  "wrong pin | insufficient balance | receiver not found | self transfer | validation error"
// @Failure      401   {object}  dto.ResponseError  "unauthorized"
// @Failure      500   {object}  dto.ResponseError  "internal server error"
// @Router       /transactions/transfer [post]
func (c *TransactionController) CreateTransfer(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var body dto.CreateTransferRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}

	err := c.transactionService.CreateTransfer(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		if errors.Is(err, appError.SelfTransferNotAllowed) ||
			errors.Is(err, appError.ReceiverNotFound) ||
			errors.Is(err, appError.InsufficientBalance) ||
			errors.Is(err, appError.WrongPin) ||
			errors.Is(err, appError.EmptyPin) {
			response.Error(ctx, 400, err.Error())
			return
		}
		response.Error(ctx, 500, "internal server error")
		return
	}

	response.Success(ctx, 201, "Transfer successful", nil)
}
