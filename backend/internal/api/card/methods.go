package card

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	_ "github.com/masharpik/vtb-hack-back/docs"
	"github.com/masharpik/vtb-hack-back/internal/api/card/models"
	jwtProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	serviceModels "github.com/masharpik/vtb-hack-back/internal/service/card/models"
	paymentServiceModels "github.com/masharpik/vtb-hack-back/internal/service/payment/models"
	"github.com/masharpik/vtb-hack-back/pkg/httputil"
)

// List godoc
// @Summary      Получить список карт
// @Description  Возвращает список карт пользователя
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.ListResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      403 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/cards [get]
func (s *Server) List(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return httputil.JSONError(c, http.StatusUnauthorized, models.ErrEmptyToken.Error())
	}

	claims, err := s.dependencies.JWTProvider.ValidateToken(c.Request().Context(), authHeader)
	if err != nil {
		err = fmt.Errorf("s.dependencies.JWTProvider.ValidateToken: %w", err)

		switch {
		case errors.Is(err, jwtProviderModels.ErrInvalidSigning),
			errors.Is(err, jwtProviderModels.ErrExpiredToken),
			errors.Is(err, jwtProviderModels.ErrInvalidToken):
			return httputil.JSONError(c, http.StatusUnauthorized, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	cards, err := s.dependencies.Service.GetCardBalancesByLogin(c.Request().Context(),
		serviceModels.GetCardBalancesByLoginRequest{
			Login: claims.Login,
		})
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.GetCardBalancesByLogin: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	var result = &models.ListResponse{
		Cards: make([]models.Card, 0, len(cards.CardBalances)),
	}

	for _, currCard := range cards.CardBalances {
		result.Cards = append(result.Cards, models.Card{
			CardID:     currCard.CardID,
			CardNumber: currCard.CardNumber,
			CardName:   currCard.CardName,
			Balance:    currCard.AccountBalance,
			Currency:   currCard.AccountCurrency,
			Bank:       currCard.Bank,
			Type:       currCard.CardType,
		})
	}

	return httputil.JSONResponse(c, http.StatusOK, result)
}

// Get godoc
// @Summary      Получить данные по карте
// @Description  Возвращает данные по карте пользователя
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param id path string true "Номер карты" example("4081781023701077248")
// @Success      200 {object} models.GetResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      404 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/card/{id} [get]
func (s *Server) Get(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return httputil.JSONError(c, http.StatusUnauthorized, models.ErrEmptyToken.Error())
	}

	claims, err := s.dependencies.JWTProvider.ValidateToken(c.Request().Context(), authHeader)
	if err != nil {
		err = fmt.Errorf("ValidateToken: %w", err)

		switch {
		case errors.Is(err, jwtProviderModels.ErrInvalidSigning),
			errors.Is(err, jwtProviderModels.ErrExpiredToken),
			errors.Is(err, jwtProviderModels.ErrInvalidToken):
			return httputil.JSONError(c, http.StatusUnauthorized, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	cardID := c.Param("id")
	if cardID == "" {
		return httputil.JSONError(c, http.StatusBadRequest, models.ErrEmptyCardID.Error())
	}

	resultCard, err := s.dependencies.Service.GetCardBalanceByLogin(c.Request().Context(),
		serviceModels.GetCardBalanceByLoginRequest{
			Login:  claims.Login,
			CardID: cardID,
		})
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.GetCardBalanceByLogin: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		case errors.Is(err, openapiProviderModels.ErrCardNotFound),
			errors.Is(err, serviceModels.ErrCardNotFound):
			return httputil.JSONError(c, http.StatusNotFound, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	if resultCard == nil || resultCard.CardBalance.CardID == "" {
		return httputil.JSONError(c, http.StatusNotFound, models.ErrCardNotFound.Error())
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.CardData{
		Number:         resultCard.CardBalance.AccountNumber,
		CardID:         resultCard.CardBalance.CardID,
		CardNumber:     resultCard.CardBalance.CardNumber,
		CardNumberFull: resultCard.CardBalance.CardNumber,
		CardName:       resultCard.CardBalance.CardName,
		HolderName:     resultCard.CardBalance.HolderName,
		ExpiryMonth:    resultCard.CardBalance.ExpiryMonth,
		ExpiryYear:     resultCard.CardBalance.ExpiryYear,
		DailyLimit:     resultCard.CardBalance.DailyLimit,
		MonthlyLimit:   resultCard.CardBalance.MonthlyLimit,
		Balance:        resultCard.CardBalance.AccountBalance,
		Currency:       resultCard.CardBalance.AccountCurrency,
		Bank:           resultCard.CardBalance.Bank,
		Type:           resultCard.CardBalance.CardType,
	})
}

// Payment godoc
// @Summary      Перевести на другой счет
// @Description  Производит перевод между счетами
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        request body models.PaymentRequest true "Запрос на перевод"
// @Success      200 {object} models.PaymentRequest
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      404 {object} httputil.ErrorResponse
// @Failure      422 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/card/payment [post]
func (s *Server) Payment(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return httputil.JSONError(c, http.StatusUnauthorized, models.ErrEmptyToken.Error())
	}

	claims, err := s.dependencies.JWTProvider.ValidateToken(c.Request().Context(), authHeader)
	if err != nil {
		err = fmt.Errorf("ValidateToken: %w", err)

		switch {
		case errors.Is(err, jwtProviderModels.ErrInvalidSigning),
			errors.Is(err, jwtProviderModels.ErrExpiredToken),
			errors.Is(err, jwtProviderModels.ErrInvalidToken):
			return httputil.JSONError(c, http.StatusUnauthorized, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	var req models.PaymentRequest
	if err = c.Bind(&req); err != nil {
		return httputil.JSONError(c, http.StatusBadRequest, fmt.Errorf("c.Bind: %w", err).Error())
	}

	if req.CardIDFrom == "" {
		return httputil.JSONError(c, http.StatusBadRequest, models.ErrEmptyCardIDFrom.Error())
	}

	if req.CardIDTo == "" {
		return httputil.JSONError(c, http.StatusBadRequest, models.ErrEmptyCardIDTo.Error())
	}

	resultPayment, err := s.dependencies.PaymentService.Payment(c.Request().Context(),
		paymentServiceModels.PaymentRequest{
			Login:         claims.Login,
			AccountIDFrom: req.CardIDFrom,
			AccountIDTo:   req.CardIDTo,
			Amount:        req.Amount,
		})
	if err != nil {
		err = fmt.Errorf("s.dependencies.PaymentService.Payment: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		case errors.Is(err, openapiProviderModels.ErrFailedPaymentStatus),
			errors.Is(err, paymentServiceModels.ErrNotEnoughMoney):
			return httputil.JSONError(c, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, paymentServiceModels.ErrAccountFromNotFound),
			errors.Is(err, paymentServiceModels.ErrAccountToNotFound):
			return httputil.JSONError(c, http.StatusNotFound, err.Error())
		case errors.Is(err, paymentServiceModels.ErrCannotPaymentToSameAccount),
			errors.Is(err, paymentServiceModels.ErrInvalidPaymentAmount):
			return httputil.JSONError(c, http.StatusBadRequest, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.PaymentResponse{
		PaymentID: resultPayment.PaymentID,
		Status:    resultPayment.Status,
	})
}
