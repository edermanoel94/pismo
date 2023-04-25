package transaction

import (
	"errors"
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

func TestTransaction_Create(t *testing.T) {

	config.Init()

	testCases := []struct {
		desc           string
		transactionReq dto.TransactionRequest

		operationTypeExpectedIndexes operationtype.Indexed
		operationTypeExpectedErr     error

		transactionInput  domain.Transaction
		transactionOutput domain.Transaction
		expectedErr       error
	}{
		{
			"create a transaction with operation_type \"PAGAMENTO\" and amount with positive sign",
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 1,
				Amount:          123.45,
			},
			map[int]string{
				1: "pagamento",
				2: "saque",
			},
			nil,
			domain.Transaction{
				Amount:          123.45,
				AccountID:       1,
				OperationTypeID: 1,
			},
			domain.Transaction{
				Amount: 123.45,
			},
			nil,
		},
		{
			"create a transaction with operation_type \"SAQUE\" and amount with negative sign",
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 2,
				Amount:          -123.45,
			},
			map[int]string{
				1: "pagamento",
				2: "saque",
			},
			nil,
			domain.Transaction{
				Amount:          -123.45,
				AccountID:       1,
				OperationTypeID: 2,
			},
			domain.Transaction{
				Amount: -123.45,
			},
			nil,
		},
		{
			"error to create a transaction",
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 1,
				Amount:          123.45,
			},
			map[int]string{
				1: "pagamento",
			},
			nil,
			domain.Transaction{
				Amount:          123.45,
				AccountID:       1,
				OperationTypeID: 1,
			},
			domain.Transaction{},
			errors.New("error to create a transaction"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			mockTransactionRepository := new(mockTransactionRepository)
			mockOperationTypeService := new(mockOperationTypeService)

			mockTransactionRepository.On("Create", tc.transactionInput).
				Return(tc.transactionOutput, tc.expectedErr)

			mockOperationTypeService.On("FindAll").Return(tc.operationTypeExpectedIndexes, tc.operationTypeExpectedErr)

			transactionService := Transaction{
				repository:           mockTransactionRepository,
				operationTypeService: mockOperationTypeService,
			}

			transactionResp, err := transactionService.Create(tc.transactionReq)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.transactionOutput.Amount, transactionResp.Amount)
			}

			mockTransactionRepository.AssertExpectations(t)
			mockOperationTypeService.AssertExpectations(t)
		})
	}
}
