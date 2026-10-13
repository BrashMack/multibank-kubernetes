// Package models объявляет структуры для запросов по пользователям в БД и ответов из них
// и переменные для других статичных данных.
package models

type Client struct {
	ID       int64
	Login    string
	Password string
}
