package server

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	return c.validator.Struct(i)
}
