package handlers

import "github.com/alexputin/subscriptions/internal/domain"

// SubscriptionRes is the response for a subscription
type SubscriptionRes struct {
	UserID      string            `json:"user_id"`
	ServiceName string            `json:"service_name"`
	Price       int               `json:"price"`
	StartDate   domain.ShortDate  `json:"start_date" swaggertype:"string" example:"07-2025"`
	EndDate     *domain.ShortDate `json:"end_date,omitempty" swaggertype:"string" example:"07-2025"`
}

// TotalPriceRes is the response for total price
type TotalPriceRes struct {
	Total int `json:"total"`
}

// SubscriptionCreateReq is used for creating a subscription
type SubscriptionCreateReq struct {
	UserID      string            `json:"user_id" validate:"required,uuid4"`
	ServiceName string            `json:"service_name" validate:"required,min=2,max=255"`
	Price       int               `json:"price" validate:"required,min=0"`
	StartDate   domain.ShortDate  `json:"start_date" validate:"required" swaggertype:"string" example:"07-2025"`
	EndDate     *domain.ShortDate `json:"end_date,omitempty" swaggertype:"string" example:"07-2025"`
}

// SubscriptionUpdateReq is used for updating a subscription
type SubscriptionUpdateReq struct {
	Price     int               `json:"price" validate:"required,min=0"`
	StartDate domain.ShortDate  `json:"start_date" validate:"required" swaggertype:"string" example:"07-2025"`
	EndDate   *domain.ShortDate `json:"end_date,omitempty" swaggertype:"string" example:"07-2025"`
}
