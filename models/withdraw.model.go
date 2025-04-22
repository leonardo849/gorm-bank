package models

import "gorm.io/gorm"

type Withdraw struct {
	gorm.Model
	BankAccountID uint `json:"bank-account_id"`
	Amount uint `json:"amount"`
	
}