package models

import "errors"

var (
	ErrInvalidAgreementAmount  = errors.New("invalid agreement amount")
	ErrConsentIsNotApproved    = errors.New("consent is not approved")
	ErrPriorityProductNotFound = errors.New("priority product not found")
)
