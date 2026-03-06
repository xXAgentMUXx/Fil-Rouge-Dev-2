package repositories

import (
	"database/sql"
	"filrouge/interfaces/models"
)

type PaymentRepository struct {
	DB *sql.DB
}

func (r *PaymentRepository) Create(payment models.Payment) error {
	_, err := r.DB.Exec(
		`INSERT INTO payments(property_id, buyer_id, amount, stripe_session_id, status)
		 VALUES($1,$2,$3,$4,$5)`,
		payment.PropertyID,
		payment.BuyerID,
		payment.Amount,
		payment.StripeSessionID,
		payment.Status,
	)

	return err
}

func (r *PaymentRepository) UpdateStatus(sessionID string, status string) error {
	_, err := r.DB.Exec(
		`UPDATE payments SET status=$1 WHERE stripe_session_id=$2`,
		status,
		sessionID,
	)

	return err
}