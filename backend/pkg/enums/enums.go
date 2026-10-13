package enums

import (
	"strings"
	"time"
)

type Bank string

type ConsentType string

type ConsentStatus string

type AccountStatus string

type AgreementStatus string

type CardStatus string

type AccountType string

type AgreementType string

type CardType string

type BalanceType string

type PaymentStatus string

const (
	AwesomeBank Bank = "awesomebank"
	VirtualBank Bank = "virtualbank"
	SmartBank   Bank = "smartbank"
)

const (
	// Read - Доступ к счетам.
	Read ConsentType = "read"
	// Write - Доступ к переводам между счетами.
	Write ConsentType = "write"
	// Create - Доступ к продуктам.
	Create ConsentType = "create"
)

const (
	AwaitingAuthorization ConsentStatus = "awaitingauthorization"
	Revoked               ConsentStatus = "revoked"
	Rejected              ConsentStatus = "rejected"
	Authorized            ConsentStatus = "authorized"
	Approved              ConsentStatus = "approved"
	Active                ConsentStatus = "active"
)

const (
	EnabledAccountStatus AccountStatus = "enabled"
)

const (
	ActiveCardStatus CardStatus = "active"
)

const (
	ActiveAgreementStatus AgreementStatus = "active"
)

const (
	CheckingAccountType AccountType = "checking"
	CardAccountType     AccountType = "card"
	CreditAccountType   AccountType = "loan"
	SavingsAccountType  AccountType = "savings"
	DepositAccountType  AccountType = "deposit"
)

const (
	DepositAgreementType AgreementType = "deposit"
)

const (
	DebitCardType  CardType = "debit"
	CreditCardType CardType = "credit"
)

const (
	InterimAvailableBalanceType BalanceType = "interimavailable"
	InterimBookedBalanceType    BalanceType = "interimbooked"
)

const (
	PendingPaymentStatus   PaymentStatus = "acceptedsettlementpending"
	CompletedPaymentStatus PaymentStatus = "acceptedsettlementcompleted"
	FailedPaymentStatus    PaymentStatus = "acceptedsettlementdailed"
)

var (
	SupportedBanks           = []Bank{AwesomeBank, VirtualBank, SmartBank}
	SupportedConsentTypes    = []ConsentType{Read, Write, Create}
	SupportedConsentStatuses = []ConsentStatus{AwaitingAuthorization, Rejected, Authorized, Revoked, Approved, Active}
	SupportedAccountStatuses = []AccountStatus{EnabledAccountStatus}
	SupportedAccountTypes    = []AccountType{CheckingAccountType, SavingsAccountType, CardAccountType,
		CreditAccountType, DepositAccountType}
	SupportedAgreementTypes    = []AgreementType{DepositAgreementType}
	SupportedCardTypes         = []CardType{CreditCardType, DebitCardType}
	SupportedCardStatuses      = []CardStatus{ActiveCardStatus}
	SupportedAgreementStatuses = []AgreementStatus{ActiveAgreementStatus}
	SupportedBalanceTypes      = []BalanceType{InterimAvailableBalanceType, InterimBookedBalanceType}
	SupportedPaymentStatus     = []PaymentStatus{PendingPaymentStatus, CompletedPaymentStatus, FailedPaymentStatus}
)

func (t ConsentType) Permissions() []string {
	switch t {
	case Read:
		return []string{"ReadAccountsDetail", "ReadBalances", "ReadTransactionsDetail", "ReadCards"}
	default:
		return []string{}
	}
}

func (t ConsentType) Reason() string {
	switch t {
	case Read:
		return "Агрегация счетов для HackAPI"
	case Create:
		return "Финансовый агрегатор для управления продуктами"
	default:
		return ""
	}
}

func (t ConsentType) CanReadProducts() bool {
	switch t {
	case Create:
		return true
	default:
		return false
	}
}

func (t ConsentType) CanOpenProducts() bool {
	switch t {
	case Create:
		return true
	default:
		return false
	}
}

func (t ConsentType) CanCloseProducts() bool {
	return false
}

func (t ConsentType) MaxAmount() float64 {
	switch t {
	case Create:
		return 1000000.00
	default:
		return 0.00
	}
}

func (t ConsentType) ValidUntil() string {
	switch t {
	case Create:
		return time.Now().Add(365 * 24 * time.Hour).Format("2006-01-02T15:04:05")
	default:
		return time.Now().Format("2006-01-02T15:04:05")
	}
}

func (t ConsentType) ValidDuration() time.Duration {
	switch t {
	case Create:
		return 365 * 24 * time.Hour
	default:
		return time.Duration(0)
	}
}

func (b Bank) ToLowerTrimSpace() Bank {
	return Bank(strings.ToLower(strings.TrimSpace(string(b))))
}

func (t ConsentType) ToLowerTrimSpace() ConsentType {
	return ConsentType(strings.ToLower(strings.TrimSpace(string(t))))
}

func (s ConsentStatus) ToLowerTrimSpace() ConsentStatus {
	return ConsentStatus(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s AccountStatus) ToLowerTrimSpace() AccountStatus {
	return AccountStatus(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s AgreementStatus) ToLowerTrimSpace() AgreementStatus {
	return AgreementStatus(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s AccountType) ToLowerTrimSpace() AccountType {
	return AccountType(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s AgreementType) ToLowerTrimSpace() AgreementType {
	return AgreementType(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s CardType) ToLowerTrimSpace() CardType {
	return CardType(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s CardStatus) ToLowerTrimSpace() CardStatus {
	return CardStatus(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s BalanceType) ToLowerTrimSpace() BalanceType {
	return BalanceType(strings.ToLower(strings.TrimSpace(string(s))))
}

func (s PaymentStatus) ToLowerTrimSpace() PaymentStatus {
	return PaymentStatus(strings.ToLower(strings.TrimSpace(string(s))))
}
