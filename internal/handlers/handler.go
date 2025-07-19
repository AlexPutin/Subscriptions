package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexputin/subscriptions/internal/domain"
	"github.com/alexputin/subscriptions/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type subscriptionsApiHandler struct {
	service  domain.UserSubscriptionService
	validate *validator.Validate
}

func NewSubscriptionsApiHandler(service domain.UserSubscriptionService) *subscriptionsApiHandler {
	return &subscriptionsApiHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *subscriptionsApiHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/api/v1")
	group.POST("/subscriptions", h.CreateSubscription)
	group.GET("/subscriptions", h.ListSubscriptions)
	group.GET("/subscriptions/:user_id/:service_name", h.GetSubscription)
	group.PUT("/subscriptions/:user_id/:service_name", h.UpdateSubscription)
	group.DELETE("/subscriptions/:user_id/:service_name", h.DeleteSubscription)
	group.GET("/subscriptions/total", h.TotalPrice)
}

// CreateSubscription handles POST /subscriptions
func (h *subscriptionsApiHandler) CreateSubscription(c echo.Context) error {
	var req SubscriptionCreateReq
	if err := c.Bind(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err)
		return nil
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return nil
	}

	sub := domain.Subscription{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err := h.service.Create(&sub)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err)
		return nil
	}
	res := SubscriptionRes(sub)
	return c.JSON(http.StatusCreated, res)
}

// ListSubscriptions handles GET /subscriptions
func (h *subscriptionsApiHandler) ListSubscriptions(c echo.Context) error {
	userID := c.QueryParam("user_id")
	limit := 20
	offset := 0
	if l := c.QueryParam("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	subs, err := h.service.List(userID, limit, offset)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err)
		return nil
	}
	res := make([]SubscriptionRes, len(subs))

	for i, s := range subs {
		res[i] = SubscriptionRes(*s)
	}
	return c.JSON(http.StatusOK, res)
}

// GetSubscription handles GET /subscriptions/:user_id/:service_name
func (h *subscriptionsApiHandler) GetSubscription(c echo.Context) error {
	userID := c.Param("user_id")
	serviceName := c.Param("service_name")
	if userID == "" || serviceName == "" {
		utils.ResponseError(c, http.StatusBadRequest, errors.New("missing user_id or service_name"))
		return nil
	}
	sub, err := h.service.Get(userID, serviceName)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err)
		return nil
	}
	if sub == nil {
		utils.ResponseError(c, http.StatusNotFound, errors.New("subscription not found"))
		return nil
	}
	res := SubscriptionRes(*sub)
	return c.JSON(http.StatusOK, res)
}

// UpdateSubscription handles PUT /subscriptions/:user_id/:service_name
func (h *subscriptionsApiHandler) UpdateSubscription(c echo.Context) error {
	userID := c.Param("user_id")
	serviceName := c.Param("service_name")
	if userID == "" || serviceName == "" {
		utils.ResponseError(c, http.StatusBadRequest, errors.New("missing user_id or service_name"))
		return nil
	}
	var req SubscriptionUpdateReq
	if err := c.Bind(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err)
		return nil
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return nil
	}

	sub := domain.Subscription{
		UserID:      userID,
		ServiceName: serviceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     &req.StartDate,
	}

	err := h.service.Update(&sub)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err)
		return nil
	}
	res := SubscriptionRes(sub)
	return c.JSON(http.StatusOK, res)
}

// DeleteSubscription handles DELETE /subscriptions/:user_id/:service_name
func (h *subscriptionsApiHandler) DeleteSubscription(c echo.Context) error {
	userID := c.Param("user_id")
	serviceName := c.Param("service_name")
	if userID == "" || serviceName == "" {
		utils.ResponseError(c, http.StatusBadRequest, errors.New("missing user_id or service_name"))
		return nil
	}
	err := h.service.Delete(userID, serviceName)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err)
		return nil
	}
	return c.NoContent(http.StatusNoContent)
}

// TotalPrice handles GET /subscriptions/total
func (h *subscriptionsApiHandler) TotalPrice(c echo.Context) error {
	userID := c.QueryParam("user_id")
	serviceName := c.QueryParam("service_name")
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	if fromStr == "" || toStr == "" {
		utils.ResponseError(c, http.StatusBadRequest, errors.New("missing from or to date"))
		return nil
	}
	from, err := parseYearMonth(fromStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, errors.New("invalid from date format, expected YYYY-MM"))
		return nil
	}
	to, err := parseYearMonth(toStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, errors.New("invalid to date format, expected YYYY-MM"))
		return nil
	}
	total, err := h.service.TotalPrice(userID, serviceName, from, to)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err)
		return nil
	}
	res := TotalPriceRes{Total: total}
	return c.JSON(http.StatusOK, res)
}

// parseYearMonth parses a string in YYYY-MM format to time.Time (first day of month)
func parseYearMonth(s string) (t time.Time, err error) {
	return time.Parse("2006-01", s)
}
