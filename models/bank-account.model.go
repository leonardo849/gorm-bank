package models

import "gorm.io/gorm"

type BankAccount struct {
	gorm.Model
	CustomerID uint `json:"customer_id"`
	Balance uint `json:"balance"` 
}