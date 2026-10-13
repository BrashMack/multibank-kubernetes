// Package models объявляет структуры для json-запросов/ответов на http-запросы клиента
// и переменные для других статичных данных.
package models

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type ListAwaitingResponse struct {
	Banks       []enums.Bank `json:"banks" example:"SmartBank,VirtualBank"`
	CreateBanks []enums.Bank `json:"createBanks" example:"SmartBank,VirtualBank"`
}

type CreateRequest struct {
	Bank enums.Bank        `json:"bank" example:"bank"`
	Type enums.ConsentType `json:"type" example:"read"`
}

type CreateResponse struct {
	AutoApprove bool `json:"auto_approve" example:"false"`
}

type ListActiveResponse struct {
	ConfirmedBanks       []enums.Bank `json:"confirmedBanks"`
	ConfirmedCreateBanks []enums.Bank `json:"confirmedCreateBanks"`
	AwaitingBanks        []enums.Bank `json:"awaitingBanks"`
	AwaitingCreateBanks  []enums.Bank `json:"awaitingCreateBanks"`
}

type CheckRequest struct {
	Bank enums.Bank        `json:"bank" example:"bank"`
	Type enums.ConsentType `json:"type" example:"read"`
}

type CheckResponse struct {
	Status string `json:"status" example:"notgiven/awaiting/approved"`
}
