package models

import "errors"

var (
	ErrEmptyToken      = errors.New("token is empty")
	ErrEmptyAccountID  = errors.New("account id is empty")
	ErrAccountNotFound = errors.New("account not found")
)
