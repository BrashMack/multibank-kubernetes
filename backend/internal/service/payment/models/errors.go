package models

import "errors"

var (
	ErrNotEnoughMoney             = errors.New("not enough money")
	ErrAccountFromNotFound        = errors.New("account from not found")
	ErrAccountToNotFound          = errors.New("account to not found")
	ErrCannotPaymentToSameAccount = errors.New("cannot be transferred to the same account")
	ErrInvalidPaymentAmount       = errors.New("invalid payment amount")
)
