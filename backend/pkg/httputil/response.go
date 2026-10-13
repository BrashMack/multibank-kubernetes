package httputil

import (
	"github.com/labstack/echo/v4"
)

func JSONResponse(c echo.Context, code int, data any) error {
	return c.JSON(code, data)
}

func JSONError(c echo.Context, code int, errorMsg string) error {
	return c.JSON(code, ErrorResponse{
		Error: errorMsg,
	})
}
