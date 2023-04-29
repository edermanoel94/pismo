package transaction

import (
	"errors"
	"github.com/edermanoel94/pismo/internal/api/account"
	operationtype "github.com/edermanoel94/pismo/internal/api/operation_type"
	"github.com/edermanoel94/pismo/internal/api/transaction/data"
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/sirupsen/logrus"
	"math"
)

var (
	ErrNoLimitCredit = errors.New("no credit limit to perform this transaction")
)

type Transaction struct {
	repository           data.TransactionRepository
	accountService       account.Service
	operationTypeService operationtype.Service
}

func New(transactionRepository data.TransactionRepository, accountService account.Service, operationTypeService operationtype.Service) *Transaction {
	return &Transaction{
		repository:           transactionRepository,
		accountService:       accountService,
		operationTypeService: operationTypeService,
	}
}

func (a *Transaction) Create(request dto.TransactionRequest) (dto.TransactionResponse, error) {

	operationTypesIndexed, err := a.operationTypeService.FindAll()

	if err != nil {
		logrus.WithFields(logrus.Fields{"request": request}).Error(err)
		return dto.TransactionResponse{}, err
	}

	accResponse, err := a.accountService.Get(request.AccountId)

	if err != nil {
		return dto.TransactionResponse{}, err
	}

	changeAmountSignToPositive(&request)

	if operationTypesIndexed.IsSettlement(request.OperationTypeId) {

		if !canPurchase(accResponse.Balance, request.Amount) {
			return dto.TransactionResponse{}, ErrNoLimitCredit
		}

		newBalance := accResponse.Balance - request.Amount

		accResponseNewBalance, err := a.accountService.UpdateBalance(accResponse.ID, newBalance)

		if err != nil {
			logrus.WithFields(logrus.Fields{"request": request}).Error(err)
			return dto.TransactionResponse{}, err
		}

		logrus.WithFields(logrus.Fields{
			"accountId":  accResponseNewBalance.ID,
			"newBalance": accResponseNewBalance.Balance,
		}).Info("Balance updated with success")

		request.Amount = swapSign(request.Amount)
	}

	if operationTypesIndexed.IsCreditLimit(request.OperationTypeId) {
		accResponse, err := a.accountService.UpdateBalance(accResponse.ID, request.Amount)

		if err != nil {
			logrus.WithFields(logrus.Fields{"request": request}).Error(err)
			return dto.TransactionResponse{}, err
		}

		logrus.WithFields(logrus.Fields{
			"accountId": accResponse.ID,
		}).Info("New credit limit is performed")
	}

	t := domain.Transaction{
		Amount:          request.Amount,
		AccountID:       uint(request.AccountId),
		OperationTypeID: uint(request.OperationTypeId),
	}

	transaction, err := a.repository.Create(t)

	if err != nil {
		logrus.WithFields(logrus.Fields{"request": request}).Error(err)
		return dto.TransactionResponse{}, err
	}

	logrus.WithFields(logrus.Fields{
		"transactionId": transaction.ID,
		"amount":        transaction.Amount,
		"accountId":     transaction.AccountID,
	}).Info("Transaction create with success")

	return dto.TransactionResponse{
		Amount: transaction.Amount,
	}, nil
}

func changeAmountSignToPositive(request *dto.TransactionRequest) {
	if math.Signbit(request.Amount) {
		request.Amount = swapSign(request.Amount)
	}
}

func swapSign(v float64) float64 {
	return -1 * v
}

func canPurchase(balance, amount float64) bool {
	return balance >= amount
}
