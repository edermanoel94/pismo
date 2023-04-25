package data

import (
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OperationTypeRepository interface {
	List() ([]domain.OperationType, error)
}

func NewOperationTypeRepository(db *gorm.DB) *Repository {

	operationTypesMap := config.Config().GetStringMapString("operation_types")

	operationTypes := make([]*domain.OperationType, 0)

	for k := range operationTypesMap {
		operationTypes = append(operationTypes, &domain.OperationType{Description: k})
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(operationTypes)

	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *gorm.DB
}

func (a *Repository) List() ([]domain.OperationType, error) {
	var operationTypes []domain.OperationType
	if err := a.db.Find(&operationTypes).Error; err != nil {
		return nil, err
	}
	return operationTypes, nil
}
