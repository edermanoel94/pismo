package account

import (
	"errors"
	"github.com/edermanoel94/pismo/internal/api/account/dto"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockAccountRepository struct {
	mock.Mock
}

func (m *mockAccountRepository) FindById(id int) (domain.Account, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *mockAccountRepository) Create(acc domain.Account) (domain.Account, error) {
	args := m.Called(acc)
	return args.Get(0).(domain.Account), args.Error(1)
}

func TestAccount_Create(t *testing.T) {

	testCases := []struct {
		desc           string
		accountRequest dto.AccountRequest
		accountInput   domain.Account
		accountOutput  domain.Account
		expectedErr    error
	}{
		{
			"create account",
			dto.AccountRequest{
				DocumentNumber: "1238123812",
			},
			domain.Account{
				DocumentNumber: "1238123812",
			},
			domain.Account{
				DocumentNumber: "1238123812",
			},
			nil,
		},
		{
			"error to create account",
			dto.AccountRequest{
				DocumentNumber: "1238123812",
			},
			domain.Account{
				DocumentNumber: "1238123812",
			},
			domain.Account{},
			errors.New("error to create account"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			mockAccountRepository := new(mockAccountRepository)

			mockAccountRepository.On("Create", tc.accountInput).Return(tc.accountOutput, tc.expectedErr)

			accountService := Account{
				repository: mockAccountRepository,
			}

			accountResponse, err := accountService.Create(tc.accountRequest)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.accountOutput.DocumentNumber, accountResponse.DocumentNumber)
			}

			mockAccountRepository.AssertExpectations(t)
		})
	}
}

func TestAccount_Get(t *testing.T) {

	testCases := []struct {
		desc          string
		id            int
		accountInput  domain.Account
		accountOutput domain.Account
		expectedErr   error
	}{
		{
			"get account",
			1,
			domain.Account{
				DocumentNumber: "1238123812",
			},
			domain.Account{
				DocumentNumber: "1238123812",
			},
			nil,
		},
		{
			"error to get account",
			1,
			domain.Account{
				DocumentNumber: "1238123812",
			},
			domain.Account{},
			errors.New("error to get account"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			mockAccountRepository := new(mockAccountRepository)

			mockAccountRepository.On("FindById", tc.id).Return(tc.accountOutput, tc.expectedErr)

			accountService := Account{
				repository: mockAccountRepository,
			}

			accountResponse, err := accountService.Get(tc.id)

			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.accountOutput.DocumentNumber, accountResponse.DocumentNumber)
			}

			mockAccountRepository.AssertExpectations(t)
		})
	}
}
