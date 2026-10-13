package consent

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"

	_ "github.com/masharpik/vtb-hack-back/docs"
	"github.com/masharpik/vtb-hack-back/internal/api/consent/models"
	jwtProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	consentProviderModels "github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	"github.com/masharpik/vtb-hack-back/pkg/enums"
	"github.com/masharpik/vtb-hack-back/pkg/httputil"
)

// ListAwaiting godoc
// @Summary      Получить список еще не подтвержденных согласий.
// @Description  Возвращает список согласий, которые надо подтвердить в его банке.
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         consent
// @Accept       json
// @Produce      json
// @Success      200 {object} models.ListAwaitingResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/consent/awaiting [get]
func (s *Server) ListAwaiting(c echo.Context) error {
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

	banks, err := s.dependencies.Service.ListAwaiting(c.Request().Context(), consentProviderModels.ListAwaitingRequest{
		Login: claims.Login,
	})
	if err != nil {
		return httputil.JSONError(c, http.StatusInternalServerError,
			fmt.Errorf("s.dependencies.Service.ListAwaiting: %w", err).Error())
	}

	awaitingBanks := make([]enums.Bank, 0, len(banks.BankRequests))
	for _, b := range banks.BankRequests {
		awaitingBanks = append(awaitingBanks, b.Bank)
	}

	awaitingCreateBanks := make([]enums.Bank, 0, len(banks.BankCreateRequests))
	for _, b := range banks.BankCreateRequests {
		awaitingCreateBanks = append(awaitingCreateBanks, b.Bank)
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.ListAwaitingResponse{
		Banks:       awaitingBanks,
		CreateBanks: awaitingCreateBanks,
	})
}

// Create godoc
// @Summary      Подать согласие пользователя.
// @Description  Отправляет запросы в банки на подачу согласия.
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         consent
// @Accept       json
// @Produce      json
// @Param        request body models.CreateRequest true "Запрос на согласие"
// @Success      200 {object} models.CreateResponse
// @Success      202 {object} models.CreateResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      409 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/consent [post]
func (s *Server) Create(c echo.Context) error {
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

	var req models.CreateRequest
	if err = c.Bind(&req); err != nil {
		return httputil.JSONError(c, http.StatusBadRequest, fmt.Errorf("c.Bind: %w", err).Error())
	}

	req.Bank = req.Bank.ToLowerTrimSpace()

	if !slices.Contains(enums.SupportedBanks, req.Bank) {
		return httputil.JSONError(c, http.StatusBadRequest, fmt.Sprintf("unsupported bank: %s", req.Bank))
	}

	req.Type = req.Type.ToLowerTrimSpace()

	if !slices.Contains(enums.SupportedConsentTypes, req.Type) {
		return httputil.JSONError(c, http.StatusBadRequest, fmt.Sprintf("unsupported consent type: %s", req.Type))
	}

	var resp *consentProviderModels.CreateReadConsentResponse

	if resp, err = s.dependencies.Service.CreateReadConsent(c.Request().Context(),
		consentProviderModels.CreateReadConsentRequest{
			Login:       claims.Login,
			ConsentType: req.Type,
			Bank:        req.Bank,
		}); err != nil {
		err = fmt.Errorf("s.dependencies.Service.CreateConsent: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrEmptyConsentID),
			errors.Is(err, openapiProviderModels.ErrEmptyRequestID):
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		case errors.Is(err, openapiProviderModels.ErrUnknownConsentType),
			errors.Is(err, consentProviderModels.ErrUnknownConsentType):
			return httputil.JSONError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, consentProviderModels.ErrConsentAlreadyExists):
			return httputil.JSONError(c, http.StatusConflict, err.Error())
		case errors.Is(err, consentProviderModels.ErrConsentAlreadyAwaiting):
			return httputil.JSONResponse(c, http.StatusAccepted, &models.CreateResponse{
				AutoApprove: false,
			})
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.CreateResponse{
		AutoApprove: resp.AutoApproved,
	})
}

// ListActive godoc
// @Summary      Получить список активных согласий с банками.
// @Description  Возвращает список подтвержденных и ожидающих подтверждение согласий с банками.
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         consent
// @Accept       json
// @Produce      json
// @Success      200 {object} models.ListActiveResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/consent/active [get]
func (s *Server) ListActive(c echo.Context) error {
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

	banks, err := s.dependencies.Service.ListActive(c.Request().Context(), consentProviderModels.ListActiveRequest{
		Login: claims.Login,
	})
	if err != nil {
		return httputil.JSONError(c, http.StatusInternalServerError,
			fmt.Errorf("s.dependencies.Service.ListActive: %w", err).Error())
	}

	confirmedBanks := make([]enums.Bank, 0, len(banks.ConfirmedBanks))
	for _, b := range banks.ConfirmedBanks {
		confirmedBanks = append(confirmedBanks, b.Bank)
	}

	awaitingBanks := make([]enums.Bank, 0, len(banks.AwaitingBanks))
	for _, b := range banks.AwaitingBanks {
		awaitingBanks = append(awaitingBanks, b.Bank)
	}

	confirmedCreateBanks := make([]enums.Bank, 0, len(banks.ConfirmedCreateBanks))
	for _, b := range banks.ConfirmedCreateBanks {
		confirmedCreateBanks = append(confirmedCreateBanks, b.Bank)
	}

	awaitingCreateBanks := make([]enums.Bank, 0, len(banks.AwaitingCreateBanks))
	for _, b := range banks.AwaitingCreateBanks {
		awaitingCreateBanks = append(awaitingCreateBanks, b.Bank)
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.ListActiveResponse{
		ConfirmedBanks:       confirmedBanks,
		ConfirmedCreateBanks: confirmedCreateBanks,
		AwaitingBanks:        awaitingBanks,
		AwaitingCreateBanks:  awaitingCreateBanks,
	})
}
