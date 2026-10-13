package openapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
)

func (c *Client) getToken(baseUrl string) (models.TokenData, error) {
	url := fmt.Sprintf("%s/auth/bank-token?client_id=%s&client_secret=%s", baseUrl, c.teamID, c.teamSecretID)

	resp, err := c.client.Post(url, "application/json", nil)
	if err != nil {
		return models.TokenData{}, fmt.Errorf("c.client.Post: %w", err)
	}

	defer resp.Body.Close()

	var tokenResp models.GetTokenResponse

	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return models.TokenData{}, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return models.TokenData{}, models.ErrEmptyToken
	}

	return models.TokenData{
		AccessToken: tokenResp.AccessToken,
		ExpiredAt:   time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}, nil
}
