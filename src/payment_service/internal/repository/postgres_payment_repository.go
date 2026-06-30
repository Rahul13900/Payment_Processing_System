package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"payment_service/internal/models"
)

type postgresPaymentRepository struct {
	db *sql.DB
}

func NewPostgresPaymentRepository(db *sql.DB) PaymentRepository {
	return &postgresPaymentRepository{db: db}
}

// Create inserts a new payment
func (r *postgresPaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (id, user_id, amount, currency, status, payment_method, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		payment.ID,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.PaymentMethod,
		payment.Description,
	).Scan(&payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// FindByID retrieves a payment by ID
func (r *postgresPaymentRepository) FindByID(ctx context.Context, id string) (*models.Payment, error) {
	query := `
		SELECT id, user_id, amount, currency, status, payment_method, 
		       provider_transaction_id, description, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	payment := &models.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.ProviderTransactionID,
		&payment.Description,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to find payment: %w", err)
	}

	return payment, nil
}

// FindByUserID retrieves all payments for a user
func (r *postgresPaymentRepository) FindByUserID(ctx context.Context, userID int, limit int, offset int) ([]*models.Payment, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM payments WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, user_id, amount, currency, status, payment_method, 
		       provider_transaction_id, description, created_at, updated_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		payment := &models.Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.PaymentMethod,
			&payment.ProviderTransactionID,
			&payment.Description,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	return payments, total, nil
}

// Update updates a payment
func (r *postgresPaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	query := `
		UPDATE payments
		SET status = $1, provider_transaction_id = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		payment.Status,
		payment.ProviderTransactionID,
		payment.ID,
	).Scan(&payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// UpdateStatus updates only the payment status
func (r *postgresPaymentRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE payments
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}
