package operationtype

import (
	"github.com/akyoto/cache"
	"github.com/edermanoel94/pismo/internal/api/operation_type/data"
	"time"
)

const (
	cacheKey = "operationTypesIndexed"
)

type OperationType struct {
	repository data.OperationTypeRepository
	cache      *cache.Cache
}

type Indexed map[int]string

func (o *OperationType) FindAll() (Indexed, error) {

	if operationTypes, ok := o.cache.Get(cacheKey); ok {
		return operationTypes.(Indexed), nil
	}

	indexes := make(map[int]string)

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
