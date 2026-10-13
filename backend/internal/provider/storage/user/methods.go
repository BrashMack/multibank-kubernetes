package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/masharpik/vtb-hack-back/internal/provider/storage/user/models"
)

// ValidateCredentials проверяет логин и хеш пароля
func (p *provider) ValidateCredentials(ctx context.Context, login, password string) error {
	const query = `
		SELECT id, login, password
		FROM clients
		WHERE login = $1
	`

	var c = new(models.Client)

	err := p.pool.QueryRow(ctx, query, login).Scan(&c.ID, &c.Login, &c.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("p.pool.QueryRow: %w", models.ErrInvalidCredentials)
		}

		return fmt.Errorf("p.pool.QueryRow: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", models.ErrInvalidCredentials)
	}

	return nil
}
