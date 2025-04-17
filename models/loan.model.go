package models

import "gorm.io/gorm"

type Loan struct {
	gorm.Model
	BaseAmount uint `json:"base_amount"`
	InterestRate uint `json:"interest_rate"`
	TotalAmount uint `json:"total_amount"`
	BankAccountID uint `json:"bank-account_id"`
}