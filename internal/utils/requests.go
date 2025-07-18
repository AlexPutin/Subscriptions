package utils

import (
	"github.com/labstack/echo/v4"
)

func ResponseError(c echo.Context, status int, err error) {
	c.JSON(status, map[string]string{
		"error": err.Error(),
	})
}

func ResponseSuccess(c echo.Context, payload any) {
	c.JSON(200, payload)
}
