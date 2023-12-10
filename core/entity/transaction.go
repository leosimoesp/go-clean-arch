package entity

import "time"

type Transaction struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Date         time.Time
	ErrorMessage string
	ID           int64
	Sender       int64
	Recipient    int64
	Value        int64
}

func NewTransaction() *Transaction {
	now := time.Now().UTC()
	return &Transaction{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Transaction) IsValid() error {
	if t.Sender <= 0 {
		return InvalidSenderErr
	}
	if t.Recipient <= 0 {
		return InvalidSenderErr
	}
	if t.Value <= 0 {
		return InvalidAmountErr
	}
	if t.Date.IsZero() {
		return InvalidDateErr
	}
	return nil
}
