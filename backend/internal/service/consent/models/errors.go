package models

import "errors"

var (
	ErrConsentAlreadyExists     = errors.New("consent is already exists")
	ErrConsentAlreadyAwaiting   = errors.New("consent is already awaiting confirmation")
	ErrUnknownConsentType       = errors.New("unknown consent type")
	ErrUnsuppotedNoAutoApproved = errors.New("unsuppoted no auto-approved")
	ErrConsentNotActive         = errors.New("consent is not active")
)
