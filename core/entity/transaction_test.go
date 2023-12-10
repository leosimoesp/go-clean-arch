package entity_test

import (
	"testing"
	"time"

	"github.com/lbsti/go-clean-arch/core/entity"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_IsValid(t *testing.T) {
	t.Run("Should create a new transaction with success", transactionSuccess)
	t.Run("Should results error if date is zero", transactionDateErr)
	t.Run("Should results error if value is equal to zero", transactionValueEqualZeroErr)
	t.Run("Should results error if value is less than zero", transactionValueLessThanZeroErr)
	t.Run("Should results error if sender id is invalid", transactionSenderInvalidErr)
	t.Run("Should results error if recipient id is invalid", transactionRecipientInvalidErr)
}

func transactionSuccess(t *testing.T) {
	transaction := entity.NewTransaction()
	transaction.Date = time.Now()
	transaction.Value = int64(5000)
	transaction.Recipient = int64(12)
	transaction.Sender = int64(25)
	err := transaction.IsValid()
	assert.Nil(t, err)
}

func transactionDateErr(t *testing.T) {
	transaction := entity.NewTransaction()
	transaction.Value = int64(5000)
	transaction.Recipient = int64(12)
	transaction.Sender = int64(25)
	err := transaction.IsValid()
	assert.NotNil(t, err)
	assert.Error(t, err, entity.InvalidDateErr.Error())
}

func transactionValueEqualZeroErr(t *testing.T) {
	transaction := entity.NewTransaction()
	transaction.Date = time.Now()
	transaction.Recipient = int64(12)
	transaction.Sender = int64(25)
	err := transaction.IsValid()
	assert.NotNil(t, err)
	assert.Error(t, err, entity.InvalidDateErr.Error())
}

func transactionValueLessThanZeroErr(t *testing.T) {
	transaction := entity.NewTransaction()
	transaction.Date = time.Now()
	transaction.Recipient = int64(12)
	transaction.Value = int64(-1)
	transaction.Sender = int64(25)
	err := transaction.IsValid()
	assert.NotNil(t, err)
	assert.Error(t, err, entity.InvalidDateErr.Error())
}

func transactionSenderInvalidErr(t *testing.T) {
	transaction := entity.NewTransaction()
	transaction.Date = time.Now()
	transaction.Value = int64(5000)
	transaction.Recipient = int64(12)
	err := transaction.IsValid()
	assert.NotNil(t, err)
	assert.Error(t, err, entity.InvalidSenderErr.Error())
}

func transactionRecipientInvalidErr(t *testing.T) {
	transaction := entity.NewTransaction()
	transaction.Date = time.Now()
	transaction.Value = int64(5000)
	transaction.Sender = int64(25)
	err := transaction.IsValid()
	assert.NotNil(t, err)
	assert.Error(t, err, entity.InvalidRecipientErr.Error())
}
