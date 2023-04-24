package transport

import (
	"errors"
	"github.com/edermanoel94/pismo/internal/api/account"
	"github.com/edermanoel94/pismo/internal/api/account/dto"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type HTTP struct {
	service account.Service
}

func NewHTTP(svc account.Service, r *echo.Echo) {
	h := &HTTP{svc}
	accountGroup := r.Group("/accounts")
	accountGroup.POST("", h.Create)
	accountGroup.GET("/:id", h.View)
}

func (h *HTTP) Create(c echo.Context) error {

	var accReq dto.AccountRequest

	if err := c.Bind(&accReq); err != nil {
		return c.JSON(http.StatusBadRequest, "request body is invalid")
	}

	accResponse, err := h.service.Create(accReq)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, "duplicated key not allowed")
		}

		return c.JSON(http.StatusInternalServerError, "error to create account")
	}

	return c.JSON(http.StatusCreated, accResponse)
}

func (h *HTTP) View(c echo.Context) error {
	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "id cannot be alpha numeric")
	}

	accRes, err := h.service.Get(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "record not found")
		}
		return c.JSON(http.StatusInternalServerError, "error get account")
	}

	return c.JSON(http.StatusOK, accRes)
}
