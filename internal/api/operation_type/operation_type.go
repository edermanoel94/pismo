package operationtype

import (
	"github.com/akyoto/cache"
	"github.com/edermanoel94/pismo/internal/api/operation_type/data"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"time"
)

const (
	cacheKey = "operationTypesIndexed"
)

type OperationType struct {
	repository data.OperationTypeRepository
	cache      *cache.Cache
}

func New(repository data.OperationTypeRepository) *OperationType {
	return &OperationType{
		repository: repository,
		cache:      cache.New(1 * time.Hour),
	}
}

type Indexed map[int]string

func IsBalanceNegative(operationTypeName string) bool {
	operationTypesMap := config.Config().GetStringMapString("operation_types")
	return operationTypesMap[operationTypeName] == "-"
}

func (o *OperationType) FindAll() (Indexed, error) {

	if operationTypes, ok := o.cache.Get(cacheKey); ok {
		return operationTypes.(Indexed), nil
	}

	indexes := make(Indexed)

	result, err := o.repository.List()

	if err != nil {
		return nil, err
	}

	for _, operationType := range result {
		indexes[int(operationType.ID)] = operationType.Description
	}

	o.cache.Set(cacheKey, indexes, 1*time.Hour)

	return indexes, nil
}
