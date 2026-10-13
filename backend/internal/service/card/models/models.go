// Package models объявляет структуры для запросов/ответов к сервисной логике по картам
// и переменные для других статичных данных.
package models

import (
	"time"

	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type Account struct {
	AccountID  string
	Number     string
	SchemeName string
	Currency   string
	Bank       enums.Bank
	Type       enums.AccountType
}

type GetAccountsResponse struct {
	Accounts []Account
}

type Card struct {
	CardID   string
	Number   string
	Type     enums.CardType
	CardName string
}

type AccountBalanceData struct {
	AccountID   string
	Number      string
	Balance     string
	Currency    string
	SchemeName  string
	Bank        enums.Bank
	Type        enums.AccountType
	NumberCards []Card
}

type AccountBalance struct {
	AccountID  string
	Number     string
	Balance    string
	Currency   string
	SchemeName string
	Bank       enums.Bank
	Type       enums.AccountType
}

type CardInfoWithoutFullNumber struct {
	CardID          string
	CardNumber      string
	CardType        enums.CardType
	CardName        string
	HolderName      string
	ExpiryMonth     int
	ExpiryYear      int
	AccountID       string
	AccountNumber   string
	AccountBalance  string
	AccountCurrency string
	Status          enums.CardStatus
	DailyLimit      string
	MonthlyLimit    string
	IssuedAt        string
	Bank            enums.Bank
}

type GetAccountBalancesResponse struct {
	AccountBalances []AccountBalance
}

type GetCardBalancesResponse struct {
	CardBalances []CardInfoWithoutFullNumber
}

type GetCardBalanceResponse struct {
	CardBalance CardInfo
}

type GetAccountBalancesByLoginRequest struct {
	Login         string
	AccountID     string
	AccountNumber []string
}

type GetAccountBalanceByLoginRequest struct {
	Login     string
	AccountID string
}

type GetAccountBalanceByLoginResponse struct {
	AccountBalance AccountBalanceData
}

type GetAccountBalanceWithoutCardsResponse struct {
	AccountBalance AccountBalance
}

type GetAccountBalancesByLoginResponse struct {
	AccountBalances []AccountBalance
}

type GetCardBalancesByLoginRequest struct {
	Login  string
	CardID []string
}

type GetCardBalancesByLoginResponse struct {
	CardBalances []CardInfoWithoutFullNumber
}

type GetCardBalanceByLoginRequest struct {
	Login  string
	CardID string
}

type CardInfo struct {
	CardID          string
	CardNumber      string
	CardFullNumber  string
	CardType        enums.CardType
	CardName        string
	HolderName      string
	ExpiryMonth     int
	ExpiryYear      int
	AccountID       string
	AccountNumber   string
	AccountBalance  string
	AccountCurrency string
	Status          enums.CardStatus
	DailyLimit      string
	MonthlyLimit    string
	IssuedAt        string
	Bank            enums.Bank
}

type GetCardBalanceByLoginResponse struct {
	CardBalance CardInfo
}

type GetAllTransactionsByLoginRequest struct {
	Login     string
	Page      int
	Limit     int
	AccountID string
}

type Amount struct {
	Amount   string
	Currency string
}

type Merchant struct {
	MerchantID string
	Name       string
	MccCode    string
	Category   string
	City       string
	Country    string
	Address    string
}

type Transaction struct {
	TransactionID          string
	BookingDateTime        time.Time
	ValueDateTime          time.Time
	Amount                 Amount
	CreditDebitIndicator   string
	Reference              string
	RemittanceInformation  string
	TransactionInformation string
	CounterpartyAccount    string
	CounterpartyBank       string
	Merchant               *Merchant
}
type GetAllTransactionsByLoginResponse struct {
	Items []Transaction
}
