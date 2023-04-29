package transaction

import (
	accdto "github.com/edermanoel94/pismo/internal/api/account/dto"
	operationtype "github.com/edermanoel94/pismo/internal/api/operation_type"
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockTransactionRepository struct {
	mock.Mock
}

func (m *mockTransactionRepository) Create(transaction domain.Transaction) (domain.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(domain.Transaction), args.Error(1)
}

type mockOperationTypeService struct {
	mock.Mock
}

func (m *mockOperationTypeService) FindAll() (operationtype.Indexed, error) {
	args := m.Called()
	return args.Get(0).(operationtype.Indexed), args.Error(1)
}

type mockAccountService struct {
	mock.Mock
}

func (m *mockAccountService) Create(request accdto.AccountRequest) (accdto.AccountResponse, error) {
	args := m.Called(request)
	return args.Get(0).(accdto.AccountResponse), args.Error(1)
}

func (m *mockAccountService) Get(id int) (accdto.AccountResponse, error) {
	args := m.Called(id)
	return args.Get(0).(accdto.AccountResponse), args.Error(1)
}

func (m *mockAccountService) UpdateBalance(id int, newBalance float64) (accdto.AccountResponse, error) {
	args := m.Called(id, newBalance)
	return args.Get(0).(accdto.AccountResponse), args.Error(1)
}

func TestTransaction_Create(t *testing.T) {

	config.Init()

	t.Run("create a transaction with operation_type \"PAGAMENTO\"", func(t *testing.T) {

		mockTransactionRepository := new(mockTransactionRepository)
		mockOperationTypeService := new(mockOperationTypeService)
		mockAccountService := new(mockAccountService)

		transactionRequest := dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 1,
			Amount:          123.45,
		}

		transactionInput := domain.Transaction{
			Amount:          123.45,
			AccountID:       1,
			OperationTypeID: 1,
		}
		transactionOutput := domain.Transaction{
			Amount: 123.45,
		}

		accountResponse := accdto.AccountResponse{}

		mockTransactionRepository.On("Create", transactionInput).
			Return(transactionOutput, nil)

		mockAccountService.On("Get", 1).
			Return(accountResponse, nil)

		mockOperationTypeService.On("FindAll").Return(operationTypeExpectedIndexes(), nil)

		transactionService := Transaction{
			repository:           mockTransactionRepository,
			accountService:       mockAccountService,
			operationTypeService: mockOperationTypeService,
		}

		transactionResp, err := transactionService.Create(transactionRequest)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, transactionOutput.Amount, transactionResp.Amount)

		mockAccountService.AssertNumberOfCalls(t, "UpdateBalance", 0)

		mockTransactionRepository.AssertExpectations(t)
		mockOperationTypeService.AssertExpectations(t)
		mockAccountService.AssertExpectations(t)
	})

	t.Run("create a transaction with operation_type \"SAQUE\"", func(t *testing.T) {

		mockTransactionRepository := new(mockTransactionRepository)
		mockOperationTypeService := new(mockOperationTypeService)
		mockAccountService := new(mockAccountService)

		transactionRequest := dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 2,
			Amount:          123.45,
		}

		transactionInput := domain.Transaction{
			Amount:          -123.45,
			AccountID:       1,
			OperationTypeID: 2,
		}
		transactionOutput := domain.Transaction{
			Amount: -123.45,
		}

		accountResponseGetMock := accdto.AccountResponse{
			ID:      1,
			Balance: 1000,
		}

		accountResponseUpdateBalanceMock := accdto.AccountResponse{
			ID:      1,
			Balance: 1000 - 123.45,
		}

		mockOperationTypeService.On("FindAll").Return(operationTypeExpectedIndexes(), nil)

		mockAccountService.On("Get", 1).
			Return(accountResponseGetMock, nil)

		mockAccountService.On("UpdateBalance", 1, 876.55).
			Return(accountResponseUpdateBalanceMock, nil)

		mockTransactionRepository.On("Create", transactionInput).
			Return(transactionOutput, nil)

		transactionService := Transaction{
			repository:           mockTransactionRepository,
			accountService:       mockAccountService,
			operationTypeService: mockOperationTypeService,
		}

		transactionResp, err := transactionService.Create(transactionRequest)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, transactionOutput.Amount, transactionResp.Amount)

		mockAccountService.AssertNumberOfCalls(t, "UpdateBalance", 1)

		mockOperationTypeService.AssertExpectations(t)
		mockAccountService.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)
	})

	t.Run("create a transaction with operation_type \"LIMITE_DE_CREDITO\"", func(t *testing.T) {

		mockTransactionRepository := new(mockTransactionRepository)
		mockOperationTypeService := new(mockOperationTypeService)
		mockAccountService := new(mockAccountService)

		transactionRequest := dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 3,
			Amount:          1000,
		}

		transactionInput := domain.Transaction{
			Amount:          1000,
			AccountID:       1,
			OperationTypeID: 3,
		}
		transactionOutput := domain.Transaction{
			Amount: 1000,
		}

		accountResponseGetMock := accdto.AccountResponse{
			ID:      1,
			Balance: 1000,
		}

		accountResponseUpdateBalanceMock := accdto.AccountResponse{
			ID:      1,
			Balance: 1000,
		}

		mockOperationTypeService.On("FindAll").Return(operationTypeExpectedIndexes(), nil)

		mockAccountService.On("Get", 1).
			Return(accountResponseGetMock, nil)

		mockAccountService.On("UpdateBalance", 1, 1000.00).
			Return(accountResponseUpdateBalanceMock, nil)

		mockTransactionRepository.On("Create", transactionInput).
			Return(transactionOutput, nil)

		transactionService := Transaction{
			repository:           mockTransactionRepository,
			accountService:       mockAccountService,
			operationTypeService: mockOperationTypeService,
		}

		transactionResp, err := transactionService.Create(transactionRequest)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, transactionOutput.Amount, transactionResp.Amount)

		mockAccountService.AssertNumberOfCalls(t, "UpdateBalance", 1)

		mockOperationTypeService.AssertExpectations(t)
		mockAccountService.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)
	})
}

func operationTypeExpectedIndexes() operationtype.Indexed {
	return operationtype.Indexed{
		1: "pagamento",
		2: "saque",
		3: "limite_de_credito",
	}
}
