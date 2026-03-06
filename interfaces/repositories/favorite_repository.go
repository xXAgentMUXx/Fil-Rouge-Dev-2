package repositories

import (
	"database/sql"
	"filrouge/interfaces/models"
)

type FavoriteRepository struct {
	DB *sql.DB
}

func (r *FavoriteRepository) Add(userID, propertyID int) error {

	_, err := r.DB.Exec(
		`INSERT INTO favorites(user_id, property_id)
		 VALUES($1,$2)
		 ON CONFLICT DO NOTHING`,
		userID,
		propertyID,
	)

	return err
}

func (r *FavoriteRepository) Remove(userID, propertyID int) error {

	_, err := r.DB.Exec(
		`DELETE FROM favorites
		 WHERE user_id=$1 AND property_id=$2`,
		userID,
		propertyID,
	)

	return err
}

func (r *FavoriteRepository) GetByUser(userID int) ([]models.Property, error) {

	rows, err := r.DB.Query(`
		SELECT p.id, p.title, p.description, p.city, p.price, p.surface, p.is_sold
		FROM properties p
		JOIN favorites f ON p.id = f.property_id
		WHERE f.user_id=$1
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var properties []models.Property

	for rows.Next() {

		var p models.Property

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.City,
			&p.Price,
			&p.Surface,
			&p.IsSold,
		)

		if err != nil {
			return nil, err
		}

		properties = append(properties, p)
	}

	return properties, nil
}