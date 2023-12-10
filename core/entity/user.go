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
	balance   int64 //balance in cents
	mu        sync.Mutex
}

func NewUser(initialBalance int64) *User {
	now := time.Now().UTC()
	return &User{balance: initialBalance, CreatedAt: now, UpdatedAt: now}
}

func (u *User) CheckDebit(amount int64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.balance-amount < 0 {
		return InsufficientBalanceErr
	}
	return nil
}

func (u *User) IsValid() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.UserType == UserType("") {
		return UsertTypeRequiredErr
	}
	if u.balance < 0 {
		return BalanceLessThanZeroErr
	}
	return nil
}

func (u *User) Debit(amount int64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.UserType == UserType(Merchant) {
		return DebitNotAllowedErr
	}
	if u.balance-amount < 0 {
		return InsufficientBalanceErr
	}
	u.balance -= amount
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) Deposit(amount int64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if amount > 0 {
		u.balance += amount
		u.UpdatedAt = time.Now().UTC()
	}
	return nil
}

func (u *User) GetBalance() int64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.balance
}
