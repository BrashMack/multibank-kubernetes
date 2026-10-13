// Package models объявляет структуры для json-запросов/ответов на http-запросы клиента
// и переменные для других статичных данных.
package models

type LoginRequest struct {
	Login    string `json:"login" example:"login" validate:"required,min=3,max=64"`
	Password string `json:"password" example:"password" validate:"required,min=6,max=128"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY30DkwIiwibmFtZSI6IkpvaG4gRG91IiwiaXNTb2NpYWwiOnRydWV9.4pcPyMD09olPSyXnrXCjTwXyr4BsezdI1AVTmud2fU4"`
}
