// Package models объявляет структуры для json-запросов/ответов на http-запросы клиента
// и переменные для других статичных данных.
package models

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type ListResponse struct {
	Cards []Card `json:"cards"`
}

type GetResponse CardData

type Card struct {
	CardID          string         `json:"cardId" example:"card-sbank-003461"`
	CardNumber      string         `json:"cardNumber" example:"**** **** **** 1386"`
	CardName        string         `json:"cardName" example:"MasterCard Standard"`
	Balance         string         `json:"balance" example:"190.00"`
	Currency        string         `json:"currency" example:"USD"`
	Bank            enums.Bank     `json:"bank" example:"VirtualBank"`
	BackgroundImage string         `json:"backgroundImage" example:"http://localhost:8081/image.jpg"`
	Type            enums.CardType `json:"type" example:"debit/credit"`
}

type CardData struct {
	Number          string         `json:"number" example:"4081781023701077248"`
	CardID          string         `json:"cardId" example:"card-sbank-003461"`
	CardNumber      string         `json:"cardNumber" example:"**** **** **** 1386"`
	CardNumberFull  string         `json:"cardNumberFull" example:"1234 2343 1324 5368"`
	CardName        string         `json:"cardName" example:"MasterCard Standard"`
	HolderName      string         `json:"holderName" example:"ИВАНОВ ИВАН ИВАНОВИЧ (TEAM237)"`
	ExpiryMonth     int            `json:"expiryMonth" example:"10"`
	ExpiryYear      int            `json:"expiryYear" example:"2028"`
	DailyLimit      string         `json:"dailyLimit" example:"100000.00"`
	MonthlyLimit    string         `json:"monthlyLimit" example:"500000.00"`
	Balance         string         `json:"balance" example:"190.00"`
	Currency        string         `json:"currency" example:"USD"`
	Bank            enums.Bank     `json:"bank" example:"VirtualBank"`
	BackgroundImage string         `json:"backgroundImage" example:"http://localhost:8081/image.jpg"`
	Type            enums.CardType `json:"type" example:"debit/credit"`
}

type PaymentRequest struct {
	CardIDFrom string `json:"cardIdFrom" example:"4081781023701077248"`
	CardIDTo   string `json:"cardIdTo" example:"4081781023701077249"`
	Amount     string `json:"amount" example:"190.00"`
}

type PaymentResponse struct {
	PaymentID string              `json:"paymentId" example:"pay-001"`
	Status    enums.PaymentStatus `json:"status" example:"completed"`
}
