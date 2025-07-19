package domain

// Subscription represents a user's subscription to a service.
type Subscription struct {
	UserID      string     `json:"user_id" db:"user_id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	StartDate   ShortDate  `json:"start_date" db:"start_date"`
	EndDate     *ShortDate `json:"end_date,omitempty" db:"end_date"`
}
