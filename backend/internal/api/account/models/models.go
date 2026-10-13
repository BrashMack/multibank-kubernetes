// Package models объявляет структуры для json-запросов/ответов на http-запросы клиента
// и переменные для других статичных данных.
package models

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type ListResponse struct {
	Accounts []Account `json:"accounts"`
}

type GetResponse AccountData

type NumberCard struct {
	CardID          string         `json:"cardId" example:"card-sbank-003461"`
	Number          string         `json:"number" example:"**** **** **** 1386"`
	BackgroundImage string         `json:"backgroundImage"`
	Type            enums.CardType `json:"type" example:"debit/credit"`
	CardName        string         `json:"cardName" example:"MasterCard Standard"`
}

type Account struct {
	AccountID string            `json:"accountId" example:"acc-3785"`
	Number    string            `json:"number" example:"4081781023701077248"`
	Balance   string            `json:"balance" example:"190.00"`
	Currency  string            `json:"currency" example:"RUB"`
	Bank      enums.Bank        `json:"bank" example:"SmartBank"`
	Type      enums.AccountType `json:"type" example:"сhecking/сard/savings/loan"`
}

type AccountData struct {
	AccountID   string            `json:"accountId" example:"acc-3785"`
	Number      string            `json:"number" example:"4081781023701077248"`
	Balance     string            `json:"balance" example:"190.00"`
	Currency    string            `json:"currency" example:"RUB"`
	Bank        enums.Bank        `json:"bank" example:"SmartBank"`
	Type        enums.AccountType `json:"type" example:"сhecking/сard/savings/loan"`
	NumberCards []NumberCard      `json:"numberCards"`
}

type Transaction struct {
	TransactionID          string `json:"transactionId" example:"tx-1ae13a1706b3"`
	Amount                 string `json:"amount" example:"190.00"`
	Currency               string `json:"currency" example:"RUB"`
	Type                   string `json:"type" example:"credit/debit"`
	CreatedAt              string `json:"createdAt" example:"2025-11-06T03:48:32.985125Z"`
	TransactionInformation string `json:"transactionInformation" example:"Перевод с 40881463ea2bad54a8"`
}

type ListTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}
