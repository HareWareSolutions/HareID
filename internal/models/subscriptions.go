package models

import (
	"HareID/internal/enums/subscription"
	"time"
)

// ID, PriceID, Status, CurrentPeriodEnd
type Subscription struct {
	ID               uint64                    `json:"id,omitempty"`
	SubscriptionID   string                    `json:"subscription_id,omitempty"`
	PriceID          string                    `json:"price_id,omitempty"`
	Status           subscription.Subscription `json:"status,omitempty"`
	CurrentPeriodEnd time.Time                 `json:"current_period_end,omitempty"`
}
