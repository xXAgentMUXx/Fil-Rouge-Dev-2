package models

import "time"

type Sale struct {
	ID         int
	PropertyID int
	BuyerID    int
	SalePrice  float64
	SoldAt     time.Time
}