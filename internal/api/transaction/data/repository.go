package data

import (
	"github.com/edermanoel94/pismo/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	Create(transaction domain.Transaction) (domain.Transaction, error)
}

func NewTransactionRepository(db *gorm.DB) *Repository {

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

func (a *Repository) Create(tx domain.Transaction) (domain.Transaction, error) {

	result := a.db.Create(&tx)

	if result.Error != nil {
		return domain.Transaction{}, result.Error
	}

	return tx, nil
}
