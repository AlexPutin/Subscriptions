package handlers

import "time"

// SubscriptionRes is the response for a subscription
type SubscriptionRes struct {
	UserID      string     `json:"user_id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

// TotalPriceRes is the response for total price
type TotalPriceRes struct {
	Total int `json:"total"`
}

// SubscriptionCreateReq is used for creating a subscription
type SubscriptionCreateReq struct {
	UserID      string     `json:"user_id" validate:"required,uuid4"`
	ServiceName string     `json:"service_name" validate:"required,min=2,max=255"`
	Price       int        `json:"price" validate:"required,min=0"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

// SubscriptionUpdateReq is used for updating a subscription
type SubscriptionUpdateReq struct {
	Price     int        `json:"price" validate:"required,min=0"`
	StartDate time.Time  `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}
