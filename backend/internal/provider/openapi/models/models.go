// Package models объявляет структуры для json-запросов/ответов открытых банковских API
// и переменные для других статичных данных.
package models

import "time"

// TokenData - структура данных токена для их использования в интерсепторе авторизации.
type TokenData struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}
