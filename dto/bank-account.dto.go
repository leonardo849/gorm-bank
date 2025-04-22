package dto



type FindBankAccountDTO struct {
	ID            uint `json:"id"`
	CustomerID    uint `json:"customer_id"`
	Balance       uint `json:"balance"`
	SentTransfers []FindBankTransferDTO `json:"sent_transfers"`
	ReceivedTransfers []FindBankTransferDTO `json:"received_transfers"`
	Loan *FindLoanDTO `json:"loan"`
	Withdrawals []FindWithdrawDTO `json:"withdrawals"`
	Deposits []FindDepositDTO `json:"deposits"`
}