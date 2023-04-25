package transaction

import (
	operationtype "github.com/edermanoel94/pismo/internal/api/operation_type"
	"github.com/edermanoel94/pismo/internal/api/transaction/data"
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/sirupsen/logrus"
	"math"
)

type Transaction struct {
	repository data.TransactionRepository

	operationTypeService operationtype.Service
}

func New(transactionRepository data.TransactionRepository, operationTypeService operationtype.Service) *Transaction {
	return &Transaction{
		repository:           transactionRepository,
		operationTypeService: operationTypeService,
	}
}

func (a *Transaction) Create(request dto.TransactionRequest) (dto.TransactionResponse, error) {

	operationTypesIndexed, err := a.operationTypeService.FindAll()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error(err)
		return dto.TransactionResponse{}, err
	}

	changeAmountSignToPositive(&request)

	t := domain.Transaction{
		Amount:          request.Amount,
		AccountID:       uint(request.AccountId),
		OperationTypeID: uint(request.OperationTypeId),
	}

	if operationtype.IsBalanceNegative(operationTypesIndexed[request.OperationTypeId]) {
		t.Amount = -1 * request.Amount
	}

	transaction, err := a.repository.Create(t)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error(err)
		return dto.TransactionResponse{}, err
	}

	return dto.TransactionResponse{
		Amount: transaction.Amount,
	}, nil
}

func changeAmountSignToPositive(request *dto.TransactionRequest) {
	if math.Signbit(request.Amount) {
		request.Amount = -1 * request.Amount
	}
}
