package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexputin/subscriptions/internal/domain"
	"github.com/alexputin/subscriptions/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	CreateFunc     func(sub *domain.Subscription) error
	GetFunc        func(userID, serviceName string) (*domain.Subscription, error)
	UpdateFunc     func(sub *domain.Subscription) error
	DeleteFunc     func(userID, serviceName string) error
	ListFunc       func(userID string, limit, offset int) ([]*domain.Subscription, error)
	TotalPriceFunc func(userID, serviceName string, from, to time.Time) (int, error)
}

func (m *mockService) Create(sub *domain.Subscription) error {
	return m.CreateFunc(sub)
}
func (m *mockService) Get(userID, serviceName string) (*domain.Subscription, error) {
	return m.GetFunc(userID, serviceName)
}
func (m *mockService) Update(sub *domain.Subscription) error {
	return m.UpdateFunc(sub)
}
func (m *mockService) Delete(userID, serviceName string) error {
	return m.DeleteFunc(userID, serviceName)
}
func (m *mockService) List(userID string, limit, offset int) ([]*domain.Subscription, error) {
	return m.ListFunc(userID, limit, offset)
}
func (m *mockService) TotalPrice(userID, serviceName string, from, to time.Time) (int, error) {
	return m.TotalPriceFunc(userID, serviceName, from, to)
}

func TestCreateSubscription(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		CreateFunc: func(sub *domain.Subscription) error {
			if sub.UserID == "fail" {
				return errors.New("fail")
			}
			return nil
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)

	body := map[string]interface{}{
		"user_id":      "user1",
		"service_name": "Netflix",
		"price":        500,
		"start_date":   "2025-07-01T00:00:00Z",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/subscriptions", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	_ = h.CreateSubscription(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetSubscription_NotFound(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		GetFunc: func(userID, serviceName string) (*domain.Subscription, error) {
			return nil, nil
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/user1/Netflix", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetParamNames("user_id", "service_name")
	c.SetParamValues("user1", "Netflix")

	_ = h.GetSubscription(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListSubscriptions(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		ListFunc: func(userID string, limit, offset int) ([]*domain.Subscription, error) {
			return []*domain.Subscription{{UserID: "user1", ServiceName: "Netflix", Price: 500, StartDate: time.Now()}}, nil
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions?user_id=user1", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	_ = h.ListSubscriptions(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTotalPrice(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		TotalPriceFunc: func(userID, serviceName string, from, to time.Time) (int, error) {
			return 1500, nil
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/total?user_id=user1&from=2025-01&to=2025-12", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	_ = h.TotalPrice(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateSubscription(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		UpdateFunc: func(sub *domain.Subscription) error {
			if sub.UserID == "fail" {
				return errors.New("fail")
			}
			return nil
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	body := map[string]interface{}{
		"price":      600,
		"start_date": "2025-07-01T00:00:00Z",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/subscriptions/user1/Netflix", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetParamNames("user_id", "service_name")
	c.SetParamValues("user1", "Netflix")

	_ = h.UpdateSubscription(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateSubscription_Error(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		UpdateFunc: func(sub *domain.Subscription) error {
			return errors.New("fail")
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	body := map[string]interface{}{
		"price":      600,
		"start_date": "2025-07-01T00:00:00Z",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/subscriptions/user1/Netflix", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetParamNames("user_id", "service_name")
	c.SetParamValues("user1", "Netflix")

	_ = h.UpdateSubscription(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteSubscription(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		DeleteFunc: func(userID, serviceName string) error {
			if userID == "fail" {
				return errors.New("fail")
			}
			return nil
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/subscriptions/user1/Netflix", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetParamNames("user_id", "service_name")
	c.SetParamValues("user1", "Netflix")

	_ = h.DeleteSubscription(c)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteSubscription_Error(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		DeleteFunc: func(userID, serviceName string) error {
			return errors.New("fail")
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/subscriptions/user1/Netflix", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetParamNames("user_id", "service_name")
	c.SetParamValues("user1", "Netflix")

	_ = h.DeleteSubscription(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateSubscription_ValidationError(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		CreateFunc: func(sub *domain.Subscription) error {
			return errors.New("validation failed: user_id is required")
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	body := map[string]interface{}{
		"service_name": "Netflix",
		"price":        500,
		"start_date":   "2025-07-01T00:00:00Z",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/subscriptions", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	_ = h.CreateSubscription(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSubscription_ValidationError(t *testing.T) {
	e := echo.New()
	ms := &mockService{
		UpdateFunc: func(sub *domain.Subscription) error {
			return errors.New("validation failed: price is required")
		},
	}
	h := handlers.NewSubscriptionsApiHandler(ms)
	body := map[string]interface{}{
		"start_date": "2025-07-01T00:00:00Z",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/subscriptions/user1/Netflix", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetParamNames("user_id", "service_name")
	c.SetParamValues("user1", "Netflix")

	_ = h.UpdateSubscription(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
