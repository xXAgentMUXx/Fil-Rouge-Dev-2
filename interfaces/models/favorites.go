package models

import "time"

type Favorite struct {
	ID         int
	UserID     int
	PropertyID int
	CreatedAt  time.Time
}