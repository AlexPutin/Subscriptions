package domain

import "time"

type UserSubscriptionRepository interface {
	Create(sub *Subscription) error
	Get(userID, serviceName string) (*Subscription, error)
	Update(sub *Subscription) error
	Delete(userID, serviceName string) error
	List(userID string, limit, offset int) ([]Subscription, error)
	TotalPrice(userID, serviceName string, from, to time.Time) (int, error)
}
