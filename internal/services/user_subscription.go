package services

import (
	"time"

	"github.com/alexputin/subscriptions/internal/domain"
)

type userSubscriptionService struct {
	repo domain.UserSubscriptionRepository
}

func NewUserSubscriptionService(repo domain.UserSubscriptionRepository) domain.UserSubscriptionService {
	return &userSubscriptionService{
		repo: repo,
	}
}

func (s *userSubscriptionService) Create(sub *domain.Subscription) error {
	return s.repo.Create(sub)
}

func (s *userSubscriptionService) Get(userID, serviceName string) (*domain.Subscription, error) {
	return s.repo.Get(userID, serviceName)
}

func (s *userSubscriptionService) Update(sub *domain.Subscription) error {
	return s.repo.Update(sub)
}

func (s *userSubscriptionService) Delete(userID, serviceName string) error {
	return s.repo.Delete(userID, serviceName)
}

func (s *userSubscriptionService) List(userID string, limit, offset int) ([]*domain.Subscription, error) {
	return s.repo.List(userID, limit, offset)
}

func (s *userSubscriptionService) TotalPrice(userID, serviceName string, from, to time.Time) (int, error) {
	return s.repo.TotalPrice(userID, serviceName, from, to)
}
