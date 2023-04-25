package data

import (
	"github.com/edermanoel94/pismo/internal/domain"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction domain.Transaction) (domain.Transaction, error)
}

func NewTransactionRepository(db *gorm.DB) *Repository {
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
