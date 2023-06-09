package operationtype

import (
	"errors"
	"github.com/akyoto/cache"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type mockOperationTypeRepository struct {
	mock.Mock
}

func (m *mockOperationTypeRepository) List() ([]domain.OperationType, error) {
	args := m.Called()
	return args.Get(0).([]domain.OperationType), args.Error(1)
}

func TestOperationType_FindAll(t *testing.T) {

	t.Run("get indexes operation types from repository", func(t *testing.T) {

		operationTypes := []domain.OperationType{
			{ID: 1, Description: "SAQUE"},
			{ID: 2, Description: "PAGAMENTO"},
		}

		mockOperationTypeRepository := new(mockOperationTypeRepository)

		mockOperationTypeRepository.On("List").Return(operationTypes, nil)

		cacheStorage := cache.New(1 * time.Hour)

		operationType := OperationType{
			cache:      cacheStorage,
			repository: mockOperationTypeRepository,
		}

		operationTypesIndexed, err := operationType.FindAll()

		assert.NoError(t, err)

		assert.Equal(t, "SAQUE", operationTypesIndexed[1])
		assert.Equal(t, "PAGAMENTO", operationTypesIndexed[2])

		assert.Equal(t, 2, len(operationTypesIndexed))

		mockOperationTypeRepository.AssertExpectations(t)
	})

	t.Run("get indexes operations types from cache", func(t *testing.T) {

		mockOperationTypeRepository := new(mockOperationTypeRepository)

		ops := []domain.OperationType{
			{ID: 1, Description: "SAQUE"},
			{ID: 2, Description: "PAGAMENTO"},
		}

		cacheStorage := cache.New(1 * time.Hour)

		// Force cache
		cacheStorage.Set(cacheKey, operationTypesIndexedFixtures(ops), 1*time.Hour)

		operationType := OperationType{
			cache:      cacheStorage,
			repository: mockOperationTypeRepository,
		}

		operationTypesIndexed, err := operationType.FindAll()

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "SAQUE", operationTypesIndexed[1])
		assert.Equal(t, "PAGAMENTO", operationTypesIndexed[2])

		assert.Equal(t, 2, len(operationTypesIndexed))

		mockOperationTypeRepository.AssertNotCalled(t, "List")

		mockOperationTypeRepository.AssertExpectations(t)
	})

	t.Run("error to get all operation types", func(t *testing.T) {

		mockOperationTypeRepository := new(mockOperationTypeRepository)

		expectedErr := errors.New("error get operation types")

		mockOperationTypeRepository.On("List").Return([]domain.OperationType{}, expectedErr)

		cacheStorage := cache.New(1 * time.Hour)

		operationType := OperationType{
			cache:      cacheStorage,
			repository: mockOperationTypeRepository,
		}

		operationTypesIndexed, err := operationType.FindAll()

		assert.Error(t, expectedErr, err)

		assert.Empty(t, operationTypesIndexed)

		mockOperationTypeRepository.AssertExpectations(t)
	})
}

func operationTypesIndexedFixtures(ops []domain.OperationType) Indexed {

	indexes := make(Indexed)

	for _, operationType := range ops {
		indexes[int(operationType.ID)] = operationType.Description
	}

	return indexes
}

func TestIsSettlement(t *testing.T) {

	config.Init()

	testCases := []struct {
		desc            string
		operationTypeId int
		expectedResult  bool
	}{
		{
			"when operationType is SAQUE",
			2,
			true,
		},
		{
			"when operationType is PAGAMENTO",
			1,
			false,
		},
		{
			"when operationType is UNKNOWN",
			1321,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			operationTypesIndexed := Indexed{
				1: "PAGAMENTO",
				2: "SAQUE",
				3: "LIMITE_DE_CREDITO",
			}

			assert.Equal(t, tc.expectedResult, operationTypesIndexed.IsSettlement(tc.operationTypeId))
		})
	}
}

func TestIndexed_IsCreditLimit(t *testing.T) {

	config.Init()

	testCases := []struct {
		desc            string
		operationTypeId int
		expectedResult  bool
	}{
		{
			"when operationType is CREDIT_LIMIT",
			3,
			true,
		},
		{
			"when operationType is SAQUE",
			2,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			operationTypesIndexed := Indexed{
				1: "PAGAMENTO",
				2: "SAQUE",
				3: "LIMITE_DE_CREDITO",
			}

			assert.Equal(t, tc.expectedResult, operationTypesIndexed.IsCreditLimit(tc.operationTypeId))
		})
	}
}
