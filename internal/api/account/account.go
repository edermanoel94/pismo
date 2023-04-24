package account

import (
	"github.com/edermanoel94/pismo/internal/api/account/data"
	"github.com/edermanoel94/pismo/internal/api/account/dto"
	"github.com/edermanoel94/pismo/internal/domain"
)

type Account struct {
	repository data.AccountRepository
}

func New(accountRepository data.AccountRepository) *Account {
	return &Account{
		repository: accountRepository,
	}
}

func (a *Account) Get(id int) (dto.AccountResponse, error) {

	acc, err := a.repository.FindById(id)

	if err != nil {
		return dto.AccountResponse{}, err
	}

	return dto.AccountResponse{
		ID:             int(acc.ID),
		DocumentNumber: acc.DocumentNumber,
	}, nil
}

func (a *Account) Create(request dto.AccountRequest) (dto.AccountResponse, error) {

	acc, err := a.repository.Create(domain.Account{
		DocumentNumber: request.DocumentNumber,
	})

	if err != nil {
		return dto.AccountResponse{}, err
	}

	return dto.AccountResponse{
		DocumentNumber: acc.DocumentNumber,
	}, nil
}
