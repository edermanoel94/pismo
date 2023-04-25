package data

import (
	"github.com/edermanoel94/pismo/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OperationTypeRepository interface {
	List() ([]domain.OperationType, error)
}

func NewOperationTypeRepository(db *gorm.DB) *Repository {

	operationTypes := []*domain.OperationType{
		{Description: "COMPRA A VISTA"},
		{Description: "COMPRA PARCELADA"},
		{Description: "SAQUE"},
		{Description: "PAGAMENTO"},
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
