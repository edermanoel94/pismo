package data

import (
	"github.com/edermanoel94/pismo/internal/domain"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindById(int) (domain.Account, error)
	Create(domain.Account) (domain.Account, error)
}

func NewAccountRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *gorm.DB
}

func (a *Repository) FindById(id int) (domain.Account, error) {
	var acc domain.Account
	acc.ID = uint(id)
	result := a.db.First(&acc, id)

	if result.Error != nil {
		return domain.Account{}, result.Error
	}

	return acc, nil
}

func (a *Repository) Create(acc domain.Account) (domain.Account, error) {

	result := a.db.Create(&acc)

	if result.Error != nil {
		return domain.Account{}, result.Error
	}

	return acc, nil
}
