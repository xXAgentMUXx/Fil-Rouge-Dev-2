package models

import "time"

type Property struct {
	ID          string
	Title       string
	Description string
	City        string
	Price       float64
	Surface     int
	AgencyID    string
	CreatedAt   time.Time
}