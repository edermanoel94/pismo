package domain

import "time"

type Account struct {
	ID             uint   `gorm:"primaryKey"`
	DocumentNumber string `gorm:"unique"`
	Transactions   []Transaction
}

type OperationType struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `gorm:"unique"`

	Transaction Transaction
}

type Transaction struct {
	ID        uint `gorm:"primaryKey"`
	Amount    float64
	EventDate time.Time `gorm:"autoCreateTime;<-:create"`

	AccountID       uint
	OperationTypeID uint
}
