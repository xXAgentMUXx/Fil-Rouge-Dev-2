package models

import "time"

type Property struct {
	ID          int
	Title       string
	Description string
	City        string
	Price       float64
	Surface     int
	AgencyID    string
	AgentID     int
	CreatedAt   time.Time
	IsSold      bool
}