package dto

import "time"

type FindWithdrawDTO struct {
	ID        uint `json:"id"`
	Amount    uint `json:"amount"`
	BankAccountID uint `json:"bank-account_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateWithdrawDTO struct {
	Amount uint `json:"amount" validate:"required,min=1"`
	
}