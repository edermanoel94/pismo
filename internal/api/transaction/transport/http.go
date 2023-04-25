package transport

import (
	"github.com/edermanoel94/pismo/internal/api/transaction"
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTTP struct {
	service transaction.Service
}

func NewHTTP(svc transaction.Service, r *echo.Echo) {
	h := &HTTP{svc}
	transactionGroup := r.Group("/transactions")
	transactionGroup.POST("", h.Create)
}

func (h *HTTP) Create(c echo.Context) error {

	var transactionReq dto.TransactionRequest

	if err := c.Bind(&transactionReq); err != nil {
		return c.JSON(http.StatusBadRequest, "request body is invalid")
	}

	accResponse, err := h.service.Create(transactionReq)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error to create transaction")
	}

	return c.JSON(http.StatusCreated, accResponse)
}
