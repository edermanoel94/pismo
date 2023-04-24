package domain

import "time"

type Account struct {
	ID             uint   `gorm:"primaryKey"`
	DocumentNumber string `gorm:"unique"`
	Transactions   []Transaction
}

type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	Amount        float64
	OperationType OperationType
	EventDate     time.Time

	AccountID uint
}

type OperationType struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `gorm:"unique"`

	TransactionID uint
}
