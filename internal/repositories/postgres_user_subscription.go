package repositories

import (
	"fmt"
	"time"

	"github.com/alexputin/subscriptions/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresUserSubscriptionRepository struct {
	db *sqlx.DB
}

func NewPostgresUserSubscriptionRepository(db *sqlx.DB) *PostgresUserSubscriptionRepository {
	return &PostgresUserSubscriptionRepository{
		db: db,
	}
}

func (r *PostgresUserSubscriptionRepository) Create(sub *domain.Subscription) error {
	_, err := r.db.Exec(`INSERT INTO subscriptions (user_id, service_name, start_date, end_date, price) VALUES ($1, $2, $3, $4, $5)`,
		sub.UserID, sub.ServiceName, sub.StartDate, sub.EndDate, sub.Price)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (r *PostgresUserSubscriptionRepository) Get(userID, serviceName string) (*domain.Subscription, error) {
	sub := &domain.Subscription{}
	err := r.db.Get(sub, `SELECT * FROM subscriptions WHERE user_id = $1 AND service_name = $2`, userID, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub, nil
}

func (r *PostgresUserSubscriptionRepository) Update(sub *domain.Subscription) error {
	_, err := r.db.Exec(`UPDATE subscriptions SET start_date = $1, end_date = $2, price = $3 WHERE user_id = $4 AND service_name = $5`,
		sub.StartDate, sub.EndDate, sub.Price, sub.UserID, sub.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (r *PostgresUserSubscriptionRepository) Delete(userID, serviceName string) error {
	_, err := r.db.Exec(`DELETE FROM subscriptions WHERE user_id = $1 AND service_name = $2`, userID, serviceName)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

func (r *PostgresUserSubscriptionRepository) List(userID string, limit, offset int) ([]domain.Subscription, error) {
	subs := make([]domain.Subscription, 0, limit)
	err := r.db.Select(&subs, `SELECT * FROM subscriptions WHERE user_id = $1 LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subs, nil
}

func (r *PostgresUserSubscriptionRepository) TotalPrice(userID, serviceName string, from, to time.Time) (int, error) {
	var totalPrice int
	err := r.db.Get(&totalPrice, `SELECT SUM(price) FROM subscriptions WHERE user_id = $1 AND service_name = $2 AND start_date >= $3 AND end_date <= $4`,
		userID, serviceName, from, to)

	if err != nil {
		return 0, fmt.Errorf("failed to calculate total price: %w", err)
	}

	return totalPrice, nil
}
