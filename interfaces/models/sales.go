package models

import "time"

type Sale struct {
	ID         string
	PropertyID string
	BuyerID    string
	SalePrice  float64
	SoldAt     time.Time
}