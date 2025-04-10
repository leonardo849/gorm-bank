package models

import "gorm.io/gorm"

type BankTransfer struct {
	gorm.Model
	Amount uint `json:"amount"`
	SenderID uint `json:"sender_id"`
	ReceiverID uint `json:"receiver_id"`
	Sender Customer  `gorm:"foreignKey:SenderID;" json:"sender"`
	Receiver Customer `gorm:"foreignKey:ReceiverID;" json:"receiver"`
}