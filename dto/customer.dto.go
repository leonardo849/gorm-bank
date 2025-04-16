package dto

import (
	"time"
)

type CreateCustomerDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,numeric,min=4,max=6"`
}

type LoginCustomerDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,numeric,min=4,max=6"`
}

type FindCustomerDTO struct {
	ID            uint               `json:"id"`
	Name          string             `json:"name"`
	Role          int                `json:"role"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	RoleUpdatedAt time.Time          `json:"role_updated_at"`
	BankAccount   FindBankAccountDTO `json:"bank_account"`
	Deposits      []FindDepositDTO   `json:"deposits"`
	SentTransfers []FindBankTransferDTO `json:"sent_transfers"`
	ReceivedTransfers []FindBankTransferDTO `json:"received_transfers"`

}
