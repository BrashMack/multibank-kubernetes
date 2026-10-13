package user

import (
	"context"

	"github.com/masharpik/vtb-hack-back/internal/service/user/models"
)

type Dependencies struct {
	Service Service
}

// Service - интерфейс методов, необходимых серверу для запросов к сервисной логике по пользователям.
type Service interface {
	Login(ctx context.Context, request models.LoginRequest) (models.LoginResponse, error)
}
