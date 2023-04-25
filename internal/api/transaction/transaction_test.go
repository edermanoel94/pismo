package transaction

import (
	"errors"
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
	"github.com/edermanoel94/pismo/internal/domain"
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

func TestTransaction_Create(t *testing.T) {

	testCases := []struct {
		desc              string
		transactionReq    dto.TransactionRequest
		transactionInput  domain.Transaction
		transactionOutput domain.Transaction
		expectedErr       error
	}{
		{
			"create a transaction",
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			domain.Transaction{
				Amount:          123.45,
				AccountID:       1,
				OperationTypeID: 4,
			},
			domain.Transaction{
				Amount: 123.45,
			},
			nil,
		},
		{
			"error to create a transaction",
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			domain.Transaction{
				Amount:          123.45,
				AccountID:       1,
				OperationTypeID: 4,
			},
			domain.Transaction{},
			errors.New("error to create a transaction"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			mockTransactionRepository := new(mockTransactionRepository)

			mockTransactionRepository.On("Create", tc.transactionInput).
				Return(tc.transactionOutput, tc.expectedErr)

			transactionService := Transaction{
				repository: mockTransactionRepository,
			}

			transactionResp, err := transactionService.Create(tc.transactionReq)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.transactionOutput.Amount, transactionResp.Amount)
			}

			mockTransactionRepository.AssertExpectations(t)
		})
	}
}
