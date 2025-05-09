package dto

import (
	"time"

)

type CreateBankTransferDTO struct {
	Amount                uint `json:"amount" validate:"required,min=1"`
	ReceiverBankAccountID uint `json:"receiver_bank-account_id" validate:"required"`
}

type FindBankTransferDTO struct {
	ID                    uint `json:"id"`
	Amount                uint `json:"amount"`
	ReceiverBankAccountID uint `json:"receiver_bank-account_id"`
	SenderBankAccountID   uint `json:"sender_bank-account_id"`
	CreatedAt time.Time `json:"created_at"`
}
