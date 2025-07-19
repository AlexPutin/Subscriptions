package domain

import "time"

type UserSubscriptionService interface {
	Create(sub *Subscription) error
	Get(userID, serviceName string) (*Subscription, error)
	Update(sub *Subscription) error
	Delete(userID, serviceName string) error
	List(userID string, limit, offset int) ([]Subscription, error)
	// Calculate total price for a period, with optional filters
	TotalPrice(userID, serviceName string, from, to time.Time) (int, error)
}
