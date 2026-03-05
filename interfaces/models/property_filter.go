package models

type PropertyFilter struct {
	City       string
	MinPrice   *float64
	MaxPrice   *float64
	MinSurface *int
}