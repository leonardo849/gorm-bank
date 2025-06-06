package models

import (
	"banco/functionscrypto"
	"banco/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name        string      `json:"name" gorm:"<-:create"`
	Password    string      `json:"password" gorm:"<-:create"`
	BankAccount BankAccount `json:"bank_account" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;"`
	// Deposits    []Deposit   `gorm:"foreignKey:CustomerID;"`
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

func (c *Customer) BeforeDelete(tx *gorm.DB) (err error) {
	if c.Role == utils.OWNER {
		return errors.New("owner can't be deleted")
	}

	return 
}