package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	_ "github.com/masharpik/vtb-hack-back/docs"
	"github.com/masharpik/vtb-hack-back/internal/api/user/models"
	jwtProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/jwt/models"
	userProviderModels "github.com/masharpik/vtb-hack-back/internal/provider/storage/user/models"
	userServiceModels "github.com/masharpik/vtb-hack-back/internal/service/user/models"
	"github.com/masharpik/vtb-hack-back/pkg/httputil"
)

// Login godoc
// @Summary      Авторизация пользователя.
// @Description  Проверка авторизационных данных от пользователя и выдача токена авторизации пользователя.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Данные для входа"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} httputil.ErrorResponse
// @Failure      401 {object} httputil.ErrorResponse
// @Failure      500 {object} httputil.ErrorResponse
// @Router       /api/v1/login [post]
func (s *Server) Login(c echo.Context) error {
	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		return httputil.JSONError(c, http.StatusBadRequest, fmt.Errorf("c.Bind: %w", err).Error())
	}

	if err := c.Validate(&req); err != nil {
		return httputil.JSONError(c, http.StatusBadRequest, fmt.Errorf("c.Validate: %w", err).Error())
	}

	token, err := s.dependencies.Service.Login(c.Request().Context(), userServiceModels.LoginRequest{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		err = fmt.Errorf("s.dependencies.Service.Login: %w", err)

		switch {
		case errors.Is(err, jwtProviderModels.ErrEmptyLogin):
			return httputil.JSONError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, userProviderModels.ErrInvalidCredentials):
			return httputil.JSONError(c, http.StatusUnauthorized, err.Error())
		default:
			return httputil.JSONError(c, http.StatusInternalServerError, err.Error())
		}

	}

	return httputil.JSONResponse(c, http.StatusOK, &models.LoginResponse{
		Token: token.Token,
	})
}
