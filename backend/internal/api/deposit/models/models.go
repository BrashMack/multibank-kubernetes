// Package models объявляет структуры для json-запросов/ответов на http-запросы клиента
// и переменные для других статичных данных.
package models

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type Deposit struct {
	ProductID    string     `json:"productId"`
	ProductName  string     `json:"productName"`
	InterestRate string     `json:"interestRate"`
	MinAmount    string     `json:"minAmount"`
	MaxAmount    string     `json:"maxAmount"`
	Currency     string     `json:"currency"`
	TermMonths   int        `json:"termMonths"`
	Bank         enums.Bank `json:"bank"`
}

type GetPriorityResponse Deposit

type Agreement struct {
	AgreementID   string     `json:"agreementId"`
	ProductID     string     `json:"productId"`
	Name          string     `json:"name"`
	Balance       string     `json:"balance"`
	Currency      string     `json:"currency"`
	InterestRate  string     `json:"interestRate"`
	AccountNumber string     `json:"accountNumber"`
	Bank          enums.Bank `json:"bank"`
}

type ListAgreementsResponse struct {
	Deposits []Agreement `json:"deposits"`
}

type ListDepositsResponse struct {
	Deposits []Deposit `json:"deposits"`
}

type CreateAgreementRequest struct {
	ProductID       string     `json:"productId"`
	Amount          string     `json:"amount"`
	Bank            enums.Bank `json:"bank"`
	TermMonths      int        `json:"termMonths"`
	SourceAccountID string     `json:"sourceAccountId"`
}

type CreateAgreementResponse Agreement
