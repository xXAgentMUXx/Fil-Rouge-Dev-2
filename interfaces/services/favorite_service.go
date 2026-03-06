package services

import (
	"filrouge/interfaces/models"
	"filrouge/interfaces/repositories"
)

type FavoriteService struct {
	Repo *repositories.FavoriteRepository
}

func (s *FavoriteService) AddFavorite(userID, propertyID int) error {
	return s.Repo.Add(userID, propertyID)
}

func (s *FavoriteService) RemoveFavorite(userID, propertyID int) error {
	return s.Repo.Remove(userID, propertyID)
}

func (s *FavoriteService) GetFavorites(userID int) ([]models.Property, error) {
	return s.Repo.GetByUser(userID)
}