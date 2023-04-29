package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/edermanoel94/pismo/internal/api/transaction"
	"github.com/edermanoel94/pismo/internal/api/transaction/dto"
	"github.com/edermanoel94/pismo/internal/infra/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockTransactionService struct {
	mock.Mock
}

func (m *mockTransactionService) Create(request dto.TransactionRequest) (dto.TransactionResponse, error) {
	args := m.Called(request)
	return args.Get(0).(dto.TransactionResponse), args.Error(1)
}

func TestHTTP_Create(t *testing.T) {

	testCases := []struct {
		desc                    string
		requestBody             string
		transactionReq          dto.TransactionRequest
		expectedErr             error
		expectedTransactionResp dto.TransactionResponse
		expectedHttpStatus      int
	}{
		{
			"success",
			`{"account_id":1,"operation_type_id":4,"amount":123.45}`,
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			nil,
			dto.TransactionResponse{},
			201,
		},
		{
			"account have no balance for this transaction",
			`{"account_id":1,"operation_type_id":4,"amount":123.45}`,
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			transaction.ErrNoLimitCredit,
			dto.TransactionResponse{},
			422,
		},
		{
			"error to create transaction",
			`{"account_id":1,"operation_type_id":4,"amount":123.45}`,
			dto.TransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			errors.New("error to create transaction"),
			dto.TransactionResponse{},
			500,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			s := server.New()

			mockTransactionService := new(mockTransactionService)

			mockTransactionService.On("Create", tc.transactionReq).Return(tc.expectedTransactionResp, tc.expectedErr)

			NewHTTP(mockTransactionService, s)

			ts := httptest.NewServer(s)

			defer ts.Close()

			path := ts.URL + "/transactions"

			res, err := http.Post(path, "application/json", bytes.NewBufferString(tc.requestBody))

			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			assert.Equal(t, tc.expectedHttpStatus, res.StatusCode)

			if tc.expectedErr != nil {
				response := unmarshalResponse(res.Body)
				assert.Regexp(t, tc.expectedErr.Error(), response)
			} else {
				assert.Equal(t, tc.expectedTransactionResp, unmarshalTransactionResponse(res.Body))
			}

			mockTransactionService.AssertExpectations(t)
		})
	}
}

func unmarshalTransactionResponse(reader io.Reader) (accResp dto.TransactionResponse) {
	data, _ := io.ReadAll(reader)
	_ = json.Unmarshal(data, &accResp)
	return
}

func unmarshalResponse(reader io.Reader) string {
	data, _ := io.ReadAll(reader)
	return string(data)
}
