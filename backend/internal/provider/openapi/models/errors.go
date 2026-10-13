package models

import "errors"

var (
	ErrEmptyToken           = errors.New("token is empty")
	ErrEmptyConsentID       = errors.New("consent id is empty during auto-approve")
	ErrConsentIDNotFound    = errors.New("consent id not found")
	ErrEmptyAgreementID     = errors.New("agreement id is empty")
	ErrEmptyRequestID       = errors.New("request id is empty")
	ErrUnknownConsentStatus = errors.New("unknown consent status")
	ErrUnknownAccountStatus = errors.New("unknown account status")
	ErrUnknownCardStatus    = errors.New("unknown card status")
	ErrUnknownAccountType   = errors.New("unknown account type")
	ErrUnknownCardType      = errors.New("unknown card type")
	ErrConsentForbidden     = errors.New("consent forbidden")
	ErrCardNotFound         = errors.New("card not found")
	ErrAccountNotFound      = errors.New("account not found")
	ErrUnknownBalanceType   = errors.New("unknown balance type")
	ErrUnknownConsentType   = errors.New("unknown consent type")
	ErrUnknownPaymentStatus = errors.New("unknown payment status")
	ErrFailedPaymentStatus  = errors.New("error payment status")
	ErrInvalidAmount        = errors.New("invalid amount")
)
