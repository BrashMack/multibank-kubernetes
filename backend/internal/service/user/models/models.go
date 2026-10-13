// Package models объявляет структуры для запросов/ответов к сервисной логике по пользователям
// и переменные для других статичных данных.
package models

type LoginRequest struct {
	Login    string
	Password string
}

type LoginResponse struct {
	Token string
}
