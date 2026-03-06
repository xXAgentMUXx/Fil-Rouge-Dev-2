package services

import (
	"errors"
	"filrouge/interfaces/models"
	"filrouge/interfaces/repositories"
)

type SaleService struct {
	SaleRepo     *repositories.SaleRepository
	PropertyRepo *repositories.PropertyRepository
}

func (s *SaleService) Sell(propertyID int, buyerID int, price float64) error {
	if price <= 0 {
		return errors.New("invalid sale price")
	}

	sale := models.Sale{
		PropertyID: propertyID,
		BuyerID:    buyerID,
		SalePrice:  price,
	}

	return s.SaleRepo.Create(sale)
}