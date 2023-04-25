package transaction

import (
	"github.com/edermanoel94/pismo/internal/api/transaction/data"
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
	"github.com/edermanoel94/pismo/internal/domain"
)

type Transaction struct {
	repository data.TransactionRepository
}

func New(transactionRepository data.TransactionRepository) *Transaction {
	return &Transaction{
		repository: transactionRepository,
	}
}

func (a *Transaction) Create(request dto.TransactionRequest) (dto.TransactionResponse, error) {

	transaction, err := a.repository.Create(domain.Transaction{
		Amount:          request.Amount,
		AccountID:       uint(request.AccountId),
		OperationTypeID: uint(request.OperationTypeId),
	})

	if err != nil {
		return dto.TransactionResponse{}, err
	}

	return dto.TransactionResponse{
		Amount: transaction.Amount,
	}, nil
}
