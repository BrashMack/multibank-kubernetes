package models

import "errors"

var (
	ErrEmptyToken      = errors.New("token is empty")
	ErrEmptyCardID     = errors.New("card id is empty")
	ErrEmptyCardIDFrom = errors.New("card id from is empty")
	ErrEmptyCardIDTo   = errors.New("card id to is empty")
	ErrCardNotFound    = errors.New("card not found")
)
