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
// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body SubscriptionCreateReq true "Subscription to create"
// @Success 201 {object} SubscriptionRes
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/subscriptions [post]
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
// ListSubscriptions godoc
// @Summary List subscriptions
// @Description List subscriptions for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} SubscriptionRes
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/subscriptions [get]
	userID := c.QueryParam("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusBadRequest, errors.New("missing user_id"))
		return nil
	}

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
	fmt.Println("Subs:", subs)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err)
		return nil
	}
	res := make([]SubscriptionRes, len(subs))

	for i, s := range subs {
		res[i] = SubscriptionRes{
			UserID:      s.UserID,
			ServiceName: s.ServiceName,
			Price:       s.Price,
			StartDate:   s.StartDate,
			EndDate:     s.EndDate,
		}
	}
	return c.JSON(http.StatusOK, res)
}

// GetSubscription handles GET /subscriptions/:user_id/:service_name
func (h *subscriptionsApiHandler) GetSubscription(c echo.Context) error {
// GetSubscription godoc
// @Summary Get a subscription
// @Description Get a subscription by user ID and service name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param service_name path string true "Service Name"
// @Success 200 {object} SubscriptionRes
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/subscriptions/{user_id}/{service_name} [get]
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
// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Update a subscription by user ID and service name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param service_name path string true "Service Name"
// @Param subscription body SubscriptionUpdateReq true "Subscription update"
// @Success 200 {object} SubscriptionRes
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/subscriptions/{user_id}/{service_name} [put]
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
// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Delete a subscription by user ID and service name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param service_name path string true "Service Name"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/subscriptions/{user_id}/{service_name} [delete]
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
// TotalPrice godoc
// @Summary Get total price
// @Description Get total price for a user's subscriptions in a date range
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string false "Service Name"
// @Param from query string true "From date (YYYY-MM)"
// @Param to query string true "To date (YYYY-MM)"
// @Success 200 {object} TotalPriceRes
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/subscriptions/total [get]
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
