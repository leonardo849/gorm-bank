package models

import "gorm.io/gorm"

type Deposit struct {
	gorm.Model
	BankAccountID uint `json:"bank-account_id"`
	Amount uint `json:"amount"`
}