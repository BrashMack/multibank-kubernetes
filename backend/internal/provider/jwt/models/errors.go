package models

import "errors"

var (
	ErrEmptyLogin     = errors.New("login is empty")
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token expired")
	ErrInvalidSigning = errors.New("invalid signing method")
)
