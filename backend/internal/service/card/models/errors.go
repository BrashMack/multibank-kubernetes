package models

import "errors"

var (
	ErrCardNotFound    = errors.New("card not found")
	ErrAccountNotFound = errors.New("account not found")
)
