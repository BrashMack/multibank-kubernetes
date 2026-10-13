// Package models объявляет структуры для запросов/ответов к сервисной логике по соглашениям
// и переменные для других статичных данных.
package models

import (
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type CreateReadConsentRequest struct {
	Login       string
	ConsentType enums.ConsentType
	Bank        enums.Bank
}

type CreateReadConsentResponse struct {
	AutoApproved bool
}

type CreateCreateConsentRequest struct {
	Login       string
	ConsentType enums.ConsentType
	BankFrom    enums.Bank
}

type CreateCreateConsentResponse struct {
	RequestID    string
	ConsentID    string
	AutoApproved bool
}

type CheckConsentRequest struct {
	Login       string
	ConsentType enums.ConsentType
	Bank        enums.Bank
	ConsentID   string
}

type CreateWriteConsentRequest struct {
	Login         string
	ConsentType   enums.ConsentType
	BankFrom      enums.Bank
	DebtorAccount string
	Amount        string
}

type CreateWriteConsentResponse struct {
	RequestID    string
	ConsentID    string
	AutoApproved bool
}

type ListAwaitingRequest struct {
	Login string
}

type BankConsent struct {
	ConsentType enums.ConsentType
	Bank        enums.Bank
	ConsentID   string
}

type BankRequest struct {
	RequestType enums.ConsentType
	Bank        enums.Bank
	RequestID   string
}

type ListAwaitingResponse struct {
	BankRequests       []BankRequest
	BankCreateRequests []BankRequest
}

type ListConfirmedRequest struct {
	Login string
}

type ListConfirmedResponse struct {
	BankConsents       []BankConsent
	BankCreateConsents []BankConsent
}

type ListActiveRequest struct {
	Login string
}

type ListActiveResponse struct {
	ConfirmedBanks       []BankConsent
	ConfirmedCreateBanks []BankConsent
	AwaitingBanks        []BankRequest
	AwaitingCreateBanks  []BankRequest
}
