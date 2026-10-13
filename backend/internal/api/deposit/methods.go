package deposit

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	_ "github.com/masharpik/vtb-hack-back/docs"
	"github.com/masharpik/vtb-hack-back/internal/api/deposit/models"
	jwtProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	consentServiceModels "github.com/masharpik/vtb-hack-back/internal/service/consent/models"
	depositServiceModels "github.com/masharpik/vtb-hack-back/internal/service/deposit/models"
	"github.com/masharpik/vtb-hack-back/pkg/httputil"
)

// GetPriority godoc
// @Summary      Получить приоритетный вклад банка партнера
// @Description  Возвращает вклад банка партнера
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         deposits
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.GetPriorityResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/deposit/priority [get]
func (s *Server) GetPriority(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return httputil.JSONError(c, http.StatusUnauthorized, models.ErrEmptyToken.Error())
	}

	_, err := s.dependencies.JWTProvider.ValidateToken(c.Request().Context(), authHeader)
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

	priorityProduct, err := s.dependencies.Service.GetPriorityProduct(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.GetPriorityProduct: %w", err)

		switch {
		case errors.Is(err, depositServiceModels.ErrPriorityProductNotFound):
			return httputil.JSONError(c, http.StatusNotFound, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.GetPriorityResponse{
		ProductID:    priorityProduct.ProductID,
		ProductName:  priorityProduct.ProductName,
		InterestRate: priorityProduct.InterestRate,
		MinAmount:    priorityProduct.MinAmount,
		MaxAmount:    priorityProduct.MaxAmount,
		Currency:     priorityProduct.Currency,
		TermMonths:   priorityProduct.TermMonths,
		Bank:         priorityProduct.Bank,
	})
}

// ListAgreements godoc
// @Summary      Получить депозиты пользователя
// @Description  Возвращает депозиты пользователя
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         deposits
// @Accept       json
// @Produce      json
// @Success      200 {object} models.ListAgreementsResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      403 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/deposits/agreements [get]
func (s *Server) ListAgreements(c echo.Context) error {
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

	result, err := s.dependencies.Service.GetAgreementsByLogin(c.Request().Context(), claims.Login)
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.GetAgreementsByLogin: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden),
			errors.Is(err, consentServiceModels.ErrConsentNotActive):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		case err != nil:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	var deposits = make([]models.Agreement, 0, len(result.Agreements))

	for _, currProduct := range result.Agreements {
		deposits = append(deposits, models.Agreement{
			AgreementID:   currProduct.AgreementID,
			ProductID:     currProduct.ProductID,
			Name:          currProduct.ProductName,
			Balance:       currProduct.Amount,
			InterestRate:  currProduct.InterestRate,
			Currency:      currProduct.Currency,
			Bank:          currProduct.Bank,
			AccountNumber: currProduct.AccountNumber,
		})
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.ListAgreementsResponse{
		Deposits: deposits,
	})
}

// ListDeposits godoc
// @Summary      Получить список предложений по вкладам
// @Description  Возвращает все предложения по вкладам
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         deposits
// @Accept       json
// @Produce      json
// @Success      200 {object} models.ListDepositsResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      403 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/deposits [get]
func (s *Server) ListDeposits(c echo.Context) error {
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

	result, err := s.dependencies.Service.GetProductsByLogin(c.Request().Context(), claims.Login)
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.GetProductsByLogin: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden),
			errors.Is(err, consentServiceModels.ErrConsentNotActive):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		case err != nil:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	var deposits = make([]models.Deposit, 0, len(result.Products))

	for _, currProduct := range result.Products {
		deposits = append(deposits, models.Deposit{
			ProductID:    currProduct.ProductID,
			ProductName:  currProduct.ProductName,
			InterestRate: currProduct.InterestRate,
			MinAmount:    currProduct.MinAmount,
			MaxAmount:    currProduct.MaxAmount,
			Currency:     currProduct.Currency,
			TermMonths:   currProduct.TermMonths,
			Bank:         currProduct.Bank,
		})
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.ListDepositsResponse{
		Deposits: deposits,
	})
}

// CreateAgreement godoc
// @Summary      Создать вклад
// @Description  Создает новый договор по вкладу
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         deposits
// @Accept       json
// @Produce      json
// @Success      200 {object} models.CreateAgreementResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      403 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/deposits/agreement [post]
func (s *Server) CreateAgreement(c echo.Context) error {
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

	var req models.CreateAgreementRequest
	if err = c.Bind(&req); err != nil {
		return httputil.JSONError(c, http.StatusBadRequest, fmt.Errorf("c.Bind: %w", err).Error())
	}

	req.Bank = req.Bank.ToLowerTrimSpace()

	result, err := s.dependencies.Service.CreateAgreement(c.Request().Context(),
		depositServiceModels.CreateAgreementRequest{
			Login:           claims.Login,
			ProductID:       req.ProductID,
			Amount:          req.Amount,
			TermMonths:      req.TermMonths,
			SourceAccountID: req.SourceAccountID,
			Bank:            req.Bank,
		})
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.CreateAgreement: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden),
			errors.Is(err, depositServiceModels.ErrConsentIsNotApproved),
			errors.Is(err, consentServiceModels.ErrUnsuppotedNoAutoApproved),
			errors.Is(err, consentServiceModels.ErrConsentNotActive):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.CreateAgreementResponse{
		AgreementID:   result.AgreementID,
		ProductID:     result.ProductID,
		Name:          result.ProductName,
		Balance:       result.Amount,
		Currency:      result.Currency,
		InterestRate:  result.InterestRate,
		AccountNumber: result.AccountNumber,
		Bank:          result.Bank,
	})
}
