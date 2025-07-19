package utils

import (
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"error"`
}

func (err *ErrorResponse) Error() string {
	return err.Message
}

func ResponseError(c echo.Context, status int, err error) {
	c.JSON(status, ErrorResponse{
		Message: err.Error(),
	})
}

func ResponseSuccess(c echo.Context, payload any) {
	c.JSON(200, payload)
}
