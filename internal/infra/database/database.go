package database

import (
	"fmt"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/edermanoel94/pismo/internal/infra/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config().Get("db.host"),
		config.Config().Get("db.user"),
		config.Config().Get("db.password"),
		config.Config().Get("db.name"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&domain.Account{},
		&domain.Transaction{},
		&domain.OperationType{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
