package handlers

import (
	"github.com/alexputin/subscriptions/internal/domain"
	"github.com/labstack/echo/v4"
)

type subscriptionsApiHandler struct {
	service domain.UserSubscriptionService
}

func NewExchangeApiHandler(service domain.UserSubscriptionService) *subscriptionsApiHandler {
	return &subscriptionsApiHandler{
		service: service,
	}
}

func (h *subscriptionsApiHandler) RegisterRoutes(app *echo.Echo) {
	// group := app.Group("api/v1")
}
