package entity_test

import (
	"testing"

	"github.com/lbsti/go-clean-arch/core/entity"
	"github.com/stretchr/testify/assert"
)

func TestUser_IsValid(t *testing.T) {
	t.Run("Should result true if balance is greather than or equal to zero and type is not empty", userIsValidSuccess)
	t.Run("Should result false if balance is less than zero", userWithInvalidBalance)
	t.Run("Should result false if type is empty", userWithEmptyType)
}

func userIsValidSuccess(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)
}

func userWithInvalidBalance(t *testing.T) {
	u := entity.NewUser(int64(-1))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.NotNil(t, err)
	assert.EqualError(t, err, entity.BalanceLessThanZeroErr.Error())
}

func userWithEmptyType(t *testing.T) {
	u := entity.NewUser(int64(1))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType("")
	err := u.IsValid()
	assert.NotNil(t, err)
	assert.EqualError(t, err, entity.UsertTypeRequiredErr.Error())
}

func TestUser_CheckDebit(t *testing.T) {
	t.Run("Should result no error if user balance is greather than debit amount", userBalanceGreatherThanAmount)
	t.Run("Should result no error if user balance is equal to debit amount", userBalanceEqualToAmount)
	t.Run("Should result no error if user balance is less than debit amount", userBalanceLessThanAmount)

}

func userBalanceGreatherThanAmount(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)
	err2 := u.CheckDebit(int64(8000))
	assert.Nil(t, err2)
}

func userBalanceEqualToAmount(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)
	err2 := u.CheckDebit(int64(10000))
	assert.Nil(t, err2)
}

func userBalanceLessThanAmount(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)
	err2 := u.CheckDebit(int64(10001))
	assert.NotNil(t, err2)
	assert.EqualError(t, err2, entity.InsufficientBalanceErr.Error())
}

func TestUser_Debit(t *testing.T) {
	t.Run("Should debit a amount from user balance with success", userDebitSuccess)
	t.Run("Should return error if merchant user try to run debit", userMerchantDebitNotAllowed)
	t.Run("Should return error if user balance is insufficient", userBalanceInsufficient)
	t.Run("Should return current balance if concurrent debits happens", userConcurrentDebits)
}

func userDebitSuccess(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)
	err2 := u.Debit(int64(10000))
	assert.Nil(t, err2)
}

func userMerchantDebitNotAllowed(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Merchant)
	err := u.IsValid()
	assert.Nil(t, err)
	err2 := u.Debit(int64(10000))
	assert.NotNil(t, err2)
	assert.EqualError(t, err2, entity.DebitNotAllowedErr.Error())
}

func userBalanceInsufficient(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)
	err2 := u.Debit(int64(10001))
	assert.NotNil(t, err2)
	assert.EqualError(t, err2, entity.InsufficientBalanceErr.Error())
}

func userConcurrentDebits(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)

	chanDebits := make(chan struct{})
	defer close(chanDebits)

	for i := 0; i < 3; i++ {
		go func() {
			_ = u.Debit(int64(1000))
			chanDebits <- struct{}{}
		}()
	}

	for i := 0; i < 3; i++ {
		<-chanDebits
	}

	assert.Equal(t, int64(7000), u.GetBalance())
}

func TestUser_Deposit(t *testing.T) {
	t.Run("Should credit a amount to user common with success", userCommonDepositSuccess)
	t.Run("Should credit a amount to user merchant with success", userMerchantDepositSuccess)
}

func userCommonDepositSuccess(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Common)
	err := u.IsValid()
	assert.Nil(t, err)

	err2 := u.Deposit(int64(5000))
	assert.Nil(t, err2)
	assert.Equal(t, int64(15000), u.GetBalance())
}

func userMerchantDepositSuccess(t *testing.T) {
	u := entity.NewUser(int64(10000))
	u.Document = "12345678901"
	u.Email = "test@test.com"
	u.FullName = "Test User"
	u.UserType = entity.UserType(entity.Merchant)
	err := u.IsValid()
	assert.Nil(t, err)

	err2 := u.Deposit(int64(5000))
	assert.Nil(t, err2)
	assert.Equal(t, int64(15000), u.GetBalance())
}
