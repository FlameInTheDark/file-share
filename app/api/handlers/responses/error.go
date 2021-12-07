package responses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Error(c echo.Context, err error) error {
	return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
}