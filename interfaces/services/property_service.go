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