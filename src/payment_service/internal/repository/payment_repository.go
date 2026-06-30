package repository

import (
	"context"
	"payment_service/internal/models"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	// Create inserts a new payment
	Create(ctx context.Context, payment *models.Payment) error

	// FindByID retrieves a payment by ID
	FindByID(ctx context.Context, id string) (*models.Payment, error)

	// FindByUserID retrieves all payments for a user
	FindByUserID(ctx context.Context, userID int, limit int, offset int) ([]*models.Payment, int, error)

	// Update updates a payment
	Update(ctx context.Context, payment *models.Payment) error

	// UpdateStatus updates only the payment status
	UpdateStatus(ctx context.Context, id, status string) error
}
