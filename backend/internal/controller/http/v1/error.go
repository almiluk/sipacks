package v1

import (
	"github.com/almiluk/sipacks/internal/controller/http/v1/models"
	"github.com/labstack/echo/v4"
)

func responseWithError(c echo.Context, code int, msg string, err error) error {
	response := models.ErrorResponse{Message: msg}
	if c.Echo().Debug {
		response.Error = err.Error()
	}
	return c.JSON(code, response)
}
