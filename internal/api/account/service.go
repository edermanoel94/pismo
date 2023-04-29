package account

import (
	"github.com/edermanoel94/pismo/internal/api/account/dto"
)

type Service interface {
	Create(dto.AccountRequest) (dto.AccountResponse, error)
	Get(int) (dto.AccountResponse, error)
	UpdateBalance(int, float64) (dto.AccountResponse, error)
}
