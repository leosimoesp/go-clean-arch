package entity

import "fmt"

var (
	BalanceLessThanZeroErr = fmt.Errorf("balance must be greater than 0")
	InsufficientBalanceErr = fmt.Errorf("insufficient balance")
	UsertTypeRequiredErr   = fmt.Errorf("user type is required")
)
