package services

import (
	"filrouge/interfaces/models"
	"filrouge/interfaces/repositories"
)

type PropertyService struct {
	Repo *repositories.PropertyRepository
}

func (s *PropertyService) ListProperties() ([]models.Property, error) {
	return s.Repo.GetAll()
}

func (s *PropertyService) AddProperty(property models.Property) error {
	return s.Repo.Add(property)
}

func (s *PropertyService) Search(filter models.PropertyFilter) ([]models.Property, error) {
	return s.Repo.FindWithFilters(filter)
}

type PaymentService struct {
	Repo *repositories.PaymentRepository
}

func (s *PaymentService) CreatePayment(payment models.Payment) error {
	return s.Repo.Create(payment)
}

func (s *PaymentService) ConfirmPayment(sessionID string) error {
	return s.Repo.UpdateStatus(sessionID, "paid")
}

func (s *PropertyService) UpdateProperty(property models.Property) error {
	return s.Repo.Update(
		property.ID,
		property.Title,
		property.Description,
		property.City,
		property.Image,
	)
}

func (s *PropertyService) GetProperty(id int) (models.Property, error) {
	return s.Repo.GetByID(id)
}

func (s *PropertyService) DeleteProperty(id int) error {
	return s.Repo.Delete(id)
}

func (s *PropertyService) GetTotalSales() (float64, error) {
	return s.Repo.GetTotalSales()
}

func (s *PropertyService) GetTotalSoldProperties() (int, error) {
	return s.Repo.GetTotalSoldProperties()
}

func (s *PropertyService) GetTopCities() ([]models.CityStats, error) {
	return s.Repo.GetTopCities()
}

func (s *PropertyService) GetMostExpensive() ([]models.Property, error) {
	return s.Repo.GetMostExpensive()
}