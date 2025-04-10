package dto

type FindBankAccountDTO struct {
	ID uint `json:"id"`
	CustomerID uint `json:"customer_id"`
	Balance uint `json:"balance"`
}