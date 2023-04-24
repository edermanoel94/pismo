package account

import (
	"github.com/edermanoel94/pismo/internal/api/account/dto"
)

type Service interface {
	Create(dto.AccountRequest) (dto.AccountResponse, error)
	Get(id int) (dto.AccountResponse, error)
}
