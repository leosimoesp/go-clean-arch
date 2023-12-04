package entity

import (
	"sync"
	"time"
)

type UserType string

const (
	Common   UserType = "common"
	Merchant UserType = "merchant"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	FullName  string
	Document  string
	Email     string
	Password  string
	UserType  UserType
	ID        int64
	Balance   int64
	mu        sync.Mutex
}

func NewUser() *User {
	return &User{}
}

func (u *User) CheckDebit(amount int64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.Balance-amount < 0 {
		return InsufficientBalanceErr
	}
	return nil
}

func (u *User) IsValid() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.Balance < 0 {
		return BalanceLessThanZeroErr
	}
	if u.UserType == UserType("") {
		return UsertTypeRequiredErr
	}
	return nil
}

func (u *User) Debit(amount int64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.Balance-amount < 0 {
		return InsufficientBalanceErr
	}
	u.Balance -= amount
	return nil
}
