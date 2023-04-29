package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edermanoel94/pismo/internal/api/account/dto"
	"github.com/edermanoel94/pismo/internal/infra/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockAccountService struct {
	mock.Mock
}

func (m *mockAccountService) Get(id int) (dto.AccountResponse, error) {
	args := m.Called(id)
	return args.Get(0).(dto.AccountResponse), args.Error(1)
}

func (m *mockAccountService) Create(request dto.AccountRequest) (dto.AccountResponse, error) {
	args := m.Called(request)
	return args.Get(0).(dto.AccountResponse), args.Error(1)
}

func (m *mockAccountService) UpdateBalance(id int, newBalance float64) (dto.AccountResponse, error) {
	args := m.Called(id, newBalance)
	return args.Get(0).(dto.AccountResponse), args.Error(1)
}

func TestHTTP_View(t *testing.T) {

	testCases := []struct {
		desc                string
		expectedErr         error
		expectedAccResponse dto.AccountResponse
		expectedHttpStatus  int
	}{
		{
			"success",
			nil,
			dto.AccountResponse{
				DocumentNumber: "document_number",
			},
			200,
		},
		{
			"record not found",
			gorm.ErrRecordNotFound,
			dto.AccountResponse{},
			404,
		},
		{
			"error get account",
			errors.New("error get account"),
			dto.AccountResponse{},
			500,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			s := server.New()

			mockAccountService := new(mockAccountService)

			id := 1

			mockAccountService.On("Get", id).Return(tc.expectedAccResponse, tc.expectedErr)

			NewHTTP(mockAccountService, s)

			ts := httptest.NewServer(s)

			defer ts.Close()

			path := fmt.Sprintf("%s/accounts/%d", ts.URL, id)

			res, err := http.Get(path)

			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			assert.Equal(t, tc.expectedHttpStatus, res.StatusCode)

			if tc.expectedErr != nil {
				response := unmarshalResponse(res.Body)
				assert.Regexp(t, tc.expectedErr.Error(), response)
			} else {
				accountResponse := unmarshalAccResponse(res.Body)
				assert.Equal(t, tc.expectedAccResponse.DocumentNumber, accountResponse.DocumentNumber)
			}

			mockAccountService.AssertExpectations(t)
		})
	}
}

func TestHTTP_Create(t *testing.T) {

	testCases := []struct {
		desc                string
		requestBody         string
		accRequest          dto.AccountRequest
		expectedErr         error
		expectedAccResponse dto.AccountResponse
		expectedHttpStatus  int
	}{
		{
			"success",
			`{"document_number":"09943432212"}`,
			dto.AccountRequest{
				DocumentNumber: "09943432212",
			},
			nil,
			dto.AccountResponse{
				DocumentNumber: "08526538403",
			},
			201,
		},
		{
			"record already exists",
			`{"document_number":"09943432212"}`,
			dto.AccountRequest{
				DocumentNumber: "09943432212",
			},
			gorm.ErrDuplicatedKey,
			dto.AccountResponse{},
			409,
		},
		{
			"error to create account",
			`{"document_number":"09943432212"}`,
			dto.AccountRequest{
				DocumentNumber: "09943432212",
			},
			errors.New("error to create account"),
			dto.AccountResponse{},
			500,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			s := server.New()

			mockAccountService := new(mockAccountService)

			mockAccountService.On("Create", tc.accRequest).Return(tc.expectedAccResponse, tc.expectedErr)

			NewHTTP(mockAccountService, s)

			ts := httptest.NewServer(s)

			defer ts.Close()

			path := ts.URL + "/accounts"

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
				assert.Equal(t, tc.expectedAccResponse.DocumentNumber, unmarshalAccResponse(res.Body).DocumentNumber)
			}

			mockAccountService.AssertExpectations(t)
		})
	}
}

func unmarshalAccResponse(reader io.Reader) (accResp dto.AccountResponse) {
	data, _ := io.ReadAll(reader)
	_ = json.Unmarshal(data, &accResp)
	return
}

func unmarshalResponse(reader io.Reader) string {
	data, _ := io.ReadAll(reader)
	return string(data)
}
