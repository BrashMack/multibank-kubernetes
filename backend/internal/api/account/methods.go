package account

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	_ "github.com/masharpik/vtb-hack-back/docs"
	"github.com/masharpik/vtb-hack-back/internal/api/account/models"
	jwtProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	openapiProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/openapi/models"
	serviceModels "github.com/masharpik/vtb-hack-back/internal/service/card/models"
	"github.com/masharpik/vtb-hack-back/pkg/httputil"
)

// List godoc
// @Summary      Получить список счетов
// @Description  Возвращает список счетов пользователя
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.ListResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/accounts [get]
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

	accounts, err := s.dependencies.Service.GetAccountBalancesByLogin(c.Request().Context(),
		serviceModels.GetAccountBalancesByLoginRequest{
			Login: claims.Login,
		})
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.GetAccountBalancesByLogin: %w", err)

		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	var result = &models.ListResponse{
		Accounts: make([]models.Account, 0, len(accounts.AccountBalances)),
	}

	for _, currAccount := range accounts.AccountBalances {
		result.Accounts = append(result.Accounts, models.Account{
			AccountID: currAccount.AccountID,
			Number:    currAccount.Number,
			Balance:   currAccount.Balance,
			Currency:  currAccount.Currency,
			Bank:      currAccount.Bank,
			Type:      currAccount.Type,
		})
	}

	return httputil.JSONResponse(c, http.StatusOK, result)
}

// Get godoc
// @Summary      Получить данные по счету
// @Description  Возвращает данные по счету пользователя
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param id path string true "Номер счета" example("acc-4744")
// @Success      200 {object} models.GetResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      404 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/account/{id} [get]
func (s *Server) Get(c echo.Context) error {
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

	accountID := c.Param("id")
	if accountID == "" {
		return httputil.JSONError(c, http.StatusBadRequest, models.ErrEmptyAccountID.Error())
	}

	resultAccount, err := s.dependencies.Service.GetAccountBalanceByLogin(c.Request().Context(),
		serviceModels.GetAccountBalanceByLoginRequest{
			Login:     claims.Login,
			AccountID: accountID,
		})
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.GetAccountBalanceByLogin: %w", err)
		switch {
		case errors.Is(err, openapiProviderModels.ErrConsentForbidden):
			return httputil.JSONError(c, http.StatusForbidden, err.Error())
		case errors.Is(err, openapiProviderModels.ErrAccountNotFound),
			errors.Is(err, serviceModels.ErrAccountNotFound):
			return httputil.JSONError(c, http.StatusNotFound, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	if resultAccount == nil || resultAccount.AccountBalance.AccountID == "" {
		return httputil.JSONError(c, http.StatusNotFound, models.ErrAccountNotFound.Error())
	}

	var numberCards = make([]models.NumberCard, 0, len(resultAccount.AccountBalance.NumberCards))

	for _, currNumberCard := range resultAccount.AccountBalance.NumberCards {
		numberCards = append(numberCards, models.NumberCard{
			CardID:          currNumberCard.CardID,
			Number:          currNumberCard.Number,
			BackgroundImage: "",
			Type:            currNumberCard.Type,
			CardName:        currNumberCard.CardName,
		})
	}

	return httputil.JSONResponse(c, http.StatusOK, &models.GetResponse{
		AccountID:   resultAccount.AccountBalance.AccountID,
		Number:      resultAccount.AccountBalance.Number,
		Balance:     resultAccount.AccountBalance.Balance,
		Currency:    resultAccount.AccountBalance.Currency,
		Bank:        resultAccount.AccountBalance.Bank,
		Type:        resultAccount.AccountBalance.Type,
		NumberCards: numberCards,
	})
}

// ListTransactions godoc
// @Summary      Получить список транзакций счета
// @Description  Возвращает список транзакций счета пользователя
// @Description  Требуется заголовок Authorization: Bearer <token>
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param id path string true "Номер счета" example("acc-7248")
// @Param page query string false "Номер страницы транзакций" example("1")
// @Param limit query string false "Количество транзакций на странице" example("1")
// @Success      200 {object} models.ListTransactionsResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/account/{id}/transactions [get]
func (s *Server) ListTransactions(c echo.Context) error {
	ctx := c.Request().Context()

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return httputil.JSONError(c, http.StatusUnauthorized, models.ErrEmptyToken.Error())
	}
	claims, err := s.dependencies.JWTProvider.ValidateToken(ctx, authHeader)
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

	accountID := c.Param("id")
	if accountID == "" {
		return httputil.JSONError(c, http.StatusBadRequest, models.ErrEmptyAccountID.Error())
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	resp, err := s.dependencies.Service.GetAllTransactionsByLogin(ctx, serviceModels.GetAllTransactionsByLoginRequest{
		Login:     claims.Login,
		AccountID: accountID,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		return httputil.JSONError(c, http.StatusInternalServerError,
			fmt.Errorf("s.dependencies.Service.GetAllTransactionsByLogin: %w", err).Error())
	}

	apiTransactions := make([]models.Transaction, 0, len(resp.Items))
	for _, t := range resp.Items {
		apiTransactions = append(apiTransactions, models.Transaction{
			TransactionID:          t.TransactionID,
			Amount:                 t.Amount.Amount,
			Currency:               t.Amount.Currency,
			Type:                   strings.ToLower(t.CreditDebitIndicator),
			CreatedAt:              t.BookingDateTime.Format("2006-01-02T15:04:05.999999999Z"),
			TransactionInformation: t.TransactionInformation,
		})
	}
	return httputil.JSONResponse(c, http.StatusOK, &models.ListTransactionsResponse{Transactions: apiTransactions})
}
