package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
)

func (p *Provider) GenerateToken(ctx context.Context, login string) (string, error) {
	if login == "" {
		return "", models.ErrEmptyLogin
	}

	now := time.Now()
	claims := models.Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(p.conf.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(p.conf.Secret))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return signed, nil
}

func (p *Provider) parseToken(ctx context.Context, token string) (*models.Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &models.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrInvalidSigning
		}

		return []byte(p.conf.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt.ParseWithClaims %w: %w", models.ErrInvalidToken, err)
	}

	claims, ok := parsed.Claims.(*models.Claims)
	if !ok {
		return nil, models.ErrInvalidToken
	}

	if !parsed.Valid {
		return nil, models.ErrInvalidToken
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, models.ErrExpiredToken
	}

	if claims.Login == "" {
		return nil, models.ErrInvalidToken
	}

	return claims, nil
}

func (p *Provider) ValidateToken(ctx context.Context, token string) (*models.Claims, error) {
	tokens := strings.Split(token, "Bearer ")
	if len(tokens) != 2 {
		return nil, models.ErrInvalidToken
	}

	return p.parseToken(ctx, tokens[1])
}
