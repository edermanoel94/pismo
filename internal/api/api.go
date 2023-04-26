package api

import (
	"github.com/edermanoel94/pismo/internal/api/account"
	ad "github.com/edermanoel94/pismo/internal/api/account/data"
	at "github.com/edermanoel94/pismo/internal/api/account/transport"
	operationtype "github.com/edermanoel94/pismo/internal/api/operation_type"
	od "github.com/edermanoel94/pismo/internal/api/operation_type/data"
	"github.com/edermanoel94/pismo/internal/api/transaction"
	td "github.com/edermanoel94/pismo/internal/api/transaction/data"
	tt "github.com/edermanoel94/pismo/internal/api/transaction/transport"
	"github.com/edermanoel94/pismo/internal/infra/database"
	"github.com/edermanoel94/pismo/internal/infra/server"
)

func Start() error {

	s := server.New()

	db, err := database.New()

	if err != nil {
		return err
	}

	accountRepository := ad.NewAccountRepository(db)
	transactionRepository := td.NewTransactionRepository(db)
	operationTypeRepository := od.NewOperationTypeRepository(db)

	operationTypeService := operationtype.New(operationTypeRepository)

	accountService := account.New(accountRepository)
	transactionService := transaction.New(transactionRepository, operationTypeService)

	at.NewHTTP(accountService, s)
	tt.NewHTTP(transactionService, s)

	server.Start(s)

	return nil
}
