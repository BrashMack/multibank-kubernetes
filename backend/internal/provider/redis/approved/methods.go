package approved

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/masharpik/vtb-hack-back/internal/provider/redis/approved/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

func key(typeConsent enums.ConsentType, login string, bank enums.Bank) string {
	return fmt.Sprintf("%s:%s:%s", typeConsent, login, bank)
}

func (p *provider) SetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank,
	consentID string, ttl time.Duration) error {
	if ttl == 0 {
		ttl = p.defaultExpiration
	}

	rKey := key(typeConsent, login, bank)

	if err := p.client.Set(ctx, rKey, consentID, ttl).Err(); err != nil {
		return fmt.Errorf("p.client.HSet: %w", err)
	}

	return nil
}

func (p *provider) GetValue(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (
	string, error) {
	val, err := p.client.Get(ctx, key(typeConsent, login, bank)).Result()
	if errors.Is(err, redis.Nil) {
		if err = p.Delete(ctx, typeConsent, login, bank); err != nil {
			return "", fmt.Errorf("p.Delete: %w", err)
		}

		return "", fmt.Errorf("p.client.Get: %w", models.ErrConsentIDNotFound)
	}

	if err != nil {
		return "", fmt.Errorf("p.client.Get: %w", err)
	}

	return val, err
}

func (p *provider) Delete(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) error {
	if err := p.client.Del(ctx, key(typeConsent, login, bank)).Err(); err != nil {
		return fmt.Errorf("p.client.Del: %w", err)
	}

	return nil
}

func (p *provider) Exists(ctx context.Context, typeConsent enums.ConsentType, login string, bank enums.Bank) (
	bool, error) {
	ok, err := p.client.Exists(ctx, key(typeConsent, login, bank)).Result()
	if err != nil {
		return false, fmt.Errorf("p.client.Exists: %w", err)
	}

	return ok == 1, nil
}

func (p *provider) ExistsByLogin(ctx context.Context, typeConsent enums.ConsentType, login string) (bool, error) {
	iter := p.client.Scan(ctx, 0, fmt.Sprintf("%s:%s:*", typeConsent, login), 1).Iterator()
	for iter.Next(ctx) {
		return true, nil
	}

	if err := iter.Err(); err != nil {
		return false, fmt.Errorf("p.client.Scan: %w", err)
	}

	return false, nil
}

func (p *provider) Close() error {
	return p.client.Close()
}

func (p *provider) GetBanksByLogin(ctx context.Context, typeConsent enums.ConsentType, login string) (
	[]models.ConsentData, error) {
	pattern := fmt.Sprintf("%s:%s:*", typeConsent, login)

	var banks = make([]models.ConsentData, 0)
	iter := p.client.Scan(ctx, 0, pattern, 100).Iterator()

	for iter.Next(ctx) {
		currKey := iter.Val()
		parts := strings.Split(currKey, ":")
		if len(parts) != 3 {
			continue
		}

		val, err := p.client.Get(ctx, currKey).Result()
		if err != nil {
			// если ключ удалён между SCAN и GET — просто пропускаем
			if errors.Is(err, redis.Nil) {
				continue
			}

			return nil, fmt.Errorf("p.client.Get(%s): %w", currKey, err)
		}

		bank := parts[len(parts)-1]
		banks = append(banks, models.ConsentData{
			Bank:      enums.Bank(bank),
			ConsentID: val,
		})
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("p.client.Scan: %w", err)
	}

	return banks, nil
}
