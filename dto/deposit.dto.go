package dto

import "time"

type CreateDepositDTO struct {
	Amount uint `json:"amount" validate:"required,min=1"`
}

type FindDepositDTO struct {
	ID uint `json:"id"`
	Amount     uint      `json:"amount"`
	BankAccountID uint      `json:"bank-account_id"`
	CreatedAt  time.Time `json:"created_at"`
}
