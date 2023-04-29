package transport

import (
	"errors"
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
		return echo.NewHTTPError(http.StatusBadRequest, "request body is invalid")
	}

	if err := c.Validate(transactionReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accResponse, err := h.service.Create(transactionReq)

	if err != nil {

		if errors.Is(err, transaction.ErrNoLimitCredit) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, accResponse)
}
