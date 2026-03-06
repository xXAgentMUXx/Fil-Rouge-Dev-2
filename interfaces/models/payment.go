package models

import "time"

type Payment struct {
	ID              int
	PropertyID      int
	BuyerID         int
	Amount          float64
	StripeSessionID string
	Status          string
	CreatedAt       time.Time
}