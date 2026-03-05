package repositories

import (
	"database/sql"
	"filrouge/interfaces/models"
	"strings"
)

type PropertyRepository struct {
	DB *sql.DB
}

func (r *PropertyRepository) GetAll() ([]models.Property, error) {
	rows, err := r.DB.Query("SELECT id, title, description, city, price, surface, agency_id, created_at FROM properties")
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
			&p.AgencyID,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, p)
	}

	return properties, nil
}

func (r *PropertyRepository) Add(property models.Property) error {
	_, err := r.DB.Exec("INSERT INTO properties(title, description, city, price, surface, agency_id) VALUES($1, $2, $3, $4, $5, $6)",
		property.Title, property.Description, property.City, property.Price, property.Surface, property.AgencyID)
	return err
}

func (r *PropertyRepository) FindWithFilters(filter models.PropertyFilter) ([]models.Property, error) {
	query := `
		SELECT id, title, description, city, price, surface, agency_id
		FROM properties
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filter.City != "" {
		query += " AND city ILIKE $" + itoa(argID)
		args = append(args, "%"+filter.City+"%")
		argID++
	}

	if filter.MinPrice != nil {
		query += " AND price >= $" + itoa(argID)
		args = append(args, *filter.MinPrice)
		argID++
	}

	if filter.MaxPrice != nil {
		query += " AND price <= $" + itoa(argID)
		args = append(args, *filter.MaxPrice)
		argID++
	}

	if filter.MinSurface != nil {
		query += " AND surface >= $" + itoa(argID)
		args = append(args, *filter.MinSurface)
		argID++
	}

	rows, err := r.DB.Query(query, args...)
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
			&p.AgencyID,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, p)
	}

	return properties, nil
}

func itoa(i int) string {
	return strings.TrimSpace(string(rune('0' + i)))
}