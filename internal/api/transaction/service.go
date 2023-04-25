package transaction

import (
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
)

type Service interface {
	Create(dto.TransactionRequest) (dto.TransactionResponse, error)
}
