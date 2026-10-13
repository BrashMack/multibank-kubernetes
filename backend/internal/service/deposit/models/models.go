// Package models объявляет структуры для запросов/ответов к сервисной логике по депозиту
// и переменные для других статичных данных.
package models

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type CreateAgreementRequest struct {
	Login           string
	ProductID       string
	Amount          string
	TermMonths      int
	SourceAccountID string
	Bank            enums.Bank
}

type Agreement struct {
	AgreementID   string
	ProductID     string
	ProductName   string
	Amount        string
	AccountNumber string
	Bank          enums.Bank
	Currency      string
	InterestRate  string
}

type CreateAgreementResponse Agreement

type GetProductsRequest struct {
	Login string
	Bank  enums.Bank
}

type GetAgreementsRequest struct {
	Login string
	Bank  enums.Bank
}

type Product struct {
	ProductID    string
	ProductName  string
	Bank         enums.Bank
	Currency     string
	InterestRate string
	MinAmount    string
	MaxAmount    string
	TermMonths   int
}

type GetProductsResponse struct {
	Products []Product
}

type GetAgreementsResponse struct {
	Agreements []Agreement
}

type GetPriorityProductResponse Product
