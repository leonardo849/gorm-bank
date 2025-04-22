package dto

import "time"

type CreateLoanDTO struct {
	Amount   uint `json:"amount" validate:"required,min=100"`
	MaturityDate time.Time `json:"maturity_date" validate:"required"`
}



type FindLoanDTO struct {
	ID uint `json:"id"`
	Amount uint `json:"amount"`
	TotalAmount uint `json:"total_amount"`
	InterestRate float64 `json:"interest_rate"`
	CreatedAt time.Time `json:"created_at"`
	MaturityDate time.Time `json:"maturity_date"`
}

