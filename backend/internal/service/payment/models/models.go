// Package models объявляет структуры для запросов/ответов к сервисной логике по картам
// и переменные для других статичных данных.
package models

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type PaymentRequest struct {
	Login         string
	AccountIDFrom string
	Amount        string
	AccountIDTo   string
}

type PaymentResponse struct {
	PaymentID string
	Status    enums.PaymentStatus
}

type CardBalance struct {
	CardID     string
	SchemeName string
	Balance    string
	Currency   string
	Bank       enums.Bank
}
