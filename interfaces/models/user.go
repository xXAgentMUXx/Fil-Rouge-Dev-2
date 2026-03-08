package models

import "time"

type User struct {
	ID        string
	Email     string
	Password  *string
	Provider  string
	Role      string
	CreatedAt time.Time
}