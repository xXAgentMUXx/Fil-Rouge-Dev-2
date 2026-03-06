package repositories

import (
	"database/sql"
	"errors"
	"filrouge/interfaces/models"
)

type SaleRepository struct {
	DB *sql.DB
}

func (r *SaleRepository) Create(sale models.Sale) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Vérifier si le bien est déjà vendu
	var isSold bool
	err = tx.QueryRow(
		"SELECT is_sold FROM properties WHERE id = $1",
		sale.PropertyID,
	).Scan(&isSold)
	if err != nil {
		tx.Rollback()
		return err
	}
	if isSold {
		tx.Rollback()
		return errors.New("property already sold")
	}

	// Insérer la vente (ID auto-généré par PostgreSQL)
	_, err = tx.Exec(`
		INSERT INTO sales (property_id, buyer_id, sale_price)
		VALUES ($1, $2, $3)
	`,
		sale.PropertyID,
		sale.BuyerID,
		sale.SalePrice,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Marquer le bien comme vendu
	_, err = tx.Exec(
		"UPDATE properties SET is_sold = true WHERE id = $1",
		sale.PropertyID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}