package v1

import "github.com/labstack/echo/v4"

type errorResponse struct {
	Error string `json:"error" example:"message"`
}

func responseWithError(c echo.Context, code int, msg string) error {
	return echo.NewHTTPError(code, msg)
}
