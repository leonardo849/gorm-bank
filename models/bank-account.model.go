package models

import (

	"gorm.io/gorm"
)

type BankAccount struct {
	gorm.Model
	CustomerID uint `json:"customer_id"`
	Balance uint `json:"balance"` 
	Loan *Loan `json:"loan" gorm:"foreignKey:BankAccountID;constraint:OnDelete:CASCADE;"`
	SentTransfers []BankTransfer `json:"sent_transfers" gorm:"foreignKey:SenderBankAccountID"`
	ReceivedTransfers []BankTransfer `json:"bank_transfers" gorm:"foreignKey:ReceiverBankAccountID"`
}