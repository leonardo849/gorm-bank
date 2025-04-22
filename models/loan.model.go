package models

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	BaseAmount uint `json:"base_amount"`
	InterestRate float64 `json:"interest_rate"`
	TotalAmount uint `json:"total_amount"`
	BankAccountID uint `json:"bank-account_id"`
	MaturityDate time.Time `json:"maturity_date"`
}