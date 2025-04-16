package dto

import "gorm.io/gorm"

type CreateBankTransferDTO struct {
	gorm.Model
	Amount     uint `json:"amount" validate:"required,min=1"`
	ReceiverID uint `json:"receiver_id" validate:"required"`
}

type FindBankTransferDTO struct {
	ID uint `json:"id"`
	Amount uint `json:"amount"`
	ReceiverID uint `json:"receiver_id"`
	SenderID uint `json:"sender_id"`
}