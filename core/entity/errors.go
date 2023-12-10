package entity

import "fmt"

var (
	BalanceLessThanZeroErr = fmt.Errorf("balance must be greater than 0")
	InsufficientBalanceErr = fmt.Errorf("insufficient balance")
	UsertTypeRequiredErr   = fmt.Errorf("user type is required")
	DebitNotAllowedErr     = fmt.Errorf("user can not run debit operation")
	InvalidSenderErr       = fmt.Errorf("invalid sender")
	InvalidRecipientErr    = fmt.Errorf("invalid recipient")
	InvalidAmountErr       = fmt.Errorf("invalid amount")
	InvalidDateErr         = fmt.Errorf("invalid operation date")
)
