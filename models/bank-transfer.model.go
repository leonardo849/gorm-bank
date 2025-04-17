package models

import "gorm.io/gorm"

type BankTransfer struct {
	gorm.Model
	Amount uint `json:"amount"`
	SenderBankAccountID   uint       `json:"sender_bank_account_id"`
	SenderBankAccount     BankAccount `gorm:"foreignKey:SenderBankAccountID;references:ID" json:"sender_bank_account"`
	ReceiverBankAccountID uint       `json:"receiver_bank_account_id"`
	ReceiverBankAccount   BankAccount `gorm:"foreignKey:ReceiverBankAccountID;references:ID" json:"receiver_bank_account"`
}