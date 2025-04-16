package models

import (
	"banco/functionscrypto"
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name        string      `json:"name" gorm:"<-:create"`
	Password    string      `json:"password" gorm:"<-:create"`
	BankAccount BankAccount `json:"bank_account" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;"`
	Deposits    []Deposit   `gorm:"foreignKey:CustomerID;"`
	SentTransfers []BankTransfer `gorm:"foreignKey:SenderID"`
	ReceivedTransfers []BankTransfer `gorm:"foreignKey:ReceiverID"`
	Loan *Loan `json:"loan" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;"`
	Role int `json:"role"`
	RoleUpdatedAt time.Time `json:"role_updated_at"`
}


func (c *Customer) BeforeSave(tx *gorm.DB) (err error) {
	hashedPassword, err := functionscrypto.StringToHash(c.Password)
	if err != nil {
		return err
	}
	c.Password = string(hashedPassword)
	return nil
}