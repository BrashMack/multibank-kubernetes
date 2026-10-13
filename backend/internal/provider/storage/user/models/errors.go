package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
)
