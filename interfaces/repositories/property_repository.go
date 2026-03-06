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
	rows, err := r.DB.Query("SELECT id, title, description, city, price, surface, agency_id, is_sold, created_at FROM properties")
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
			&p.IsSold,
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
	_, err := r.DB.Exec(
	`INSERT INTO properties(title, description, city, price, surface, agency_id, agent_id)
	VALUES($1,$2,$3,$4,$5,$6,$7)`,
	property.Title,
	property.Description,
	property.City,
	property.Price,
	property.Surface,
	property.AgencyID,
	property.AgentID,
	)
	return err
}

func (r *PropertyRepository) FindWithFilters(filter models.PropertyFilter) ([]models.Property, error) {
	query := `
	SELECT id, title, description, city, price, surface, agency_id, is_sold, created_at
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
			&p.IsSold,
			&p.CreatedAt,
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

func (r *PropertyRepository) Update(property models.Property) error {

	_, err := r.DB.Exec(`
	UPDATE properties
	SET title=$1, description=$2, city=$3, price=$4, surface=$5
	WHERE id=$6`,
		property.Title,
		property.Description,
		property.City,
		property.Price,
		property.Surface,
		property.ID,
	)

	return err
}

func (r *PropertyRepository) GetByID(id int) (models.Property, error) {

	var p models.Property

	err := r.DB.QueryRow(`
	SELECT id,title,description,city,price,surface,agency_id,agent_id,is_sold,created_at
	FROM properties WHERE id=$1`, id).
	Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.City,
		&p.Price,
		&p.Surface,
		&p.AgencyID,
		&p.AgentID,
		&p.IsSold,
		&p.CreatedAt,
	)

	return p, err
}

func (r *PropertyRepository) Delete(id int) error {

	_, err := r.DB.Exec(
	"DELETE FROM properties WHERE id=$1",
	id,
	)

	return err
}

func (r *PropertyRepository) GetTotalSales() (float64, error) {

	var total float64

	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(price),0)
		FROM properties
		WHERE is_sold = true
	`).Scan(&total)

	return total, err
}

func (r *PropertyRepository) GetTotalSoldProperties() (int, error) {

	var count int

	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM properties
		WHERE is_sold = true
	`).Scan(&count)

	return count, err
}

func (r *PropertyRepository) GetTopCities() ([]models.CityStats, error) {

	rows, err := r.DB.Query(`
		SELECT city, COUNT(*) as total
		FROM properties
		GROUP BY city
		ORDER BY total DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.CityStats

	for rows.Next() {
		var c models.CityStats
		rows.Scan(&c.City, &c.Total)
		cities = append(cities, c)
	}

	return cities, nil
}

func (r *PropertyRepository) GetMostExpensive() ([]models.Property, error) {

	rows, err := r.DB.Query(`
		SELECT id, title, city, price
		FROM properties
		ORDER BY price DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []models.Property

	for rows.Next() {
		var p models.Property
		rows.Scan(&p.ID, &p.Title, &p.City, &p.Price)
		properties = append(properties, p)
	}

	return properties, nil
}