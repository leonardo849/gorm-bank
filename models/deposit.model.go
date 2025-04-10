package models

import "gorm.io/gorm"

type Deposit struct {
	gorm.Model
	CustomerID uint `json:"customer_id"`
	Amount uint `json:"amount"`
}