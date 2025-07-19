package services_test

import (
	"testing"
	"time"

	"github.com/alexputin/subscriptions/internal/domain"
	"github.com/alexputin/subscriptions/internal/services"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	CreateFunc     func(sub *domain.Subscription) error
	GetFunc        func(userID, serviceName string) (*domain.Subscription, error)
	UpdateFunc     func(sub *domain.Subscription) error
	DeleteFunc     func(userID, serviceName string) error
	ListFunc       func(userID string, limit, offset int) ([]domain.Subscription, error)
	TotalPriceFunc func(userID, serviceName string, from, to time.Time) (int, error)
}

func (m *mockRepo) Create(sub *domain.Subscription) error {
	return m.CreateFunc(sub)
}
func (m *mockRepo) Get(userID, serviceName string) (*domain.Subscription, error) {
	return m.GetFunc(userID, serviceName)
}
func (m *mockRepo) Update(sub *domain.Subscription) error {
	return m.UpdateFunc(sub)
}
func (m *mockRepo) Delete(userID, serviceName string) error {
	return m.DeleteFunc(userID, serviceName)
}
func (m *mockRepo) List(userID string, limit, offset int) ([]domain.Subscription, error) {
	return m.ListFunc(userID, limit, offset)
}
func (m *mockRepo) TotalPrice(userID, serviceName string, from, to time.Time) (int, error) {
	return m.TotalPriceFunc(userID, serviceName, from, to)
}

func TestUserSubscriptionService_Create_Ok(t *testing.T) {
	called := false
	repo := mockRepo{
		CreateFunc: func(sub *domain.Subscription) error {
			called = true
			return nil
		},
	}
	svc := services.NewUserSubscriptionService(&repo)
	sub := domain.Subscription{UserID: "user1", ServiceName: "Netflix", Price: 500, StartDate: domain.ShortDate{Time: time.Now()}}
	err := svc.Create(&sub)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestUserSubscriptionService_Update_Ok(t *testing.T) {
	called := false
	repo := mockRepo{
		UpdateFunc: func(sub *domain.Subscription) error {
			called = true
			return nil
		},
	}
	svc := services.NewUserSubscriptionService(&repo)
	sub := domain.Subscription{UserID: "user1", ServiceName: "Netflix", Price: 500, StartDate: domain.ShortDate{Time: time.Now()}}
	err := svc.Update(&sub)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestUserSubscriptionService_Get(t *testing.T) {
	repo := mockRepo{
		GetFunc: func(userID, serviceName string) (*domain.Subscription, error) {
			return &domain.Subscription{UserID: userID, ServiceName: serviceName, Price: 100}, nil
		},
	}
	svc := services.NewUserSubscriptionService(&repo)
	sub, err := svc.Get("user1", "Netflix")
	assert.NoError(t, err)
	assert.Equal(t, "user1", sub.UserID)
	assert.Equal(t, "Netflix", sub.ServiceName)
}

func TestUserSubscriptionService_Delete(t *testing.T) {
	called := false
	repo := mockRepo{
		DeleteFunc: func(userID, serviceName string) error {
			called = true
			return nil
		},
	}
	svc := services.NewUserSubscriptionService(&repo)
	err := svc.Delete("user1", "Netflix")
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestUserSubscriptionService_List(t *testing.T) {
	repo := mockRepo{
		ListFunc: func(userID string, limit, offset int) ([]domain.Subscription, error) {
			return []domain.Subscription{{UserID: userID, ServiceName: "Netflix", Price: 100}}, nil
		},
	}
	svc := services.NewUserSubscriptionService(&repo)
	list, err := svc.List("user1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "user1", list[0].UserID)
}

func TestUserSubscriptionService_TotalPrice(t *testing.T) {
	repo := mockRepo{
		TotalPriceFunc: func(userID, serviceName string, from, to time.Time) (int, error) {
			return 1234, nil
		},
	}
	svc := services.NewUserSubscriptionService(&repo)
	total, err := svc.TotalPrice("user1", "Netflix", time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 1234, total)
}
