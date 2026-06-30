package service

import (
	"context"
	"database/sql"
	"fmt"
	"payment_service/internal/models"
	"payment_service/internal/paymenterrors"
	"payment_service/internal/repository"
	"shared/logger"
	"strings"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, userID int, amount int64, currency, paymentMethod, description string) (*models.Payment, error) {
	ctx, method := logger.FuncInitializer(ctx, "CreatePayment")
	defer logger.FuncDisposer(ctx, method)

	ctx = logger.WithUserID(ctx, userID)

	// Validate inputs
	if amount <= 0 {
		logger.Warn(ctx, method, "Invalid amount: "+fmt.Sprintf("%d", amount))
		return nil, paymenterrors.ErrInvalidAmount
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))
	if !isValidCurrency(currency) {
		logger.Warn(ctx, method, "Invalid currency: "+currency)
		return nil, paymenterrors.ErrInvalidCurrency
	}

	logger.Info(ctx, method, fmt.Sprintf("Creating payment: amount=%d, currency=%s", amount, currency))

	// Generate payment ID
	paymentID := "pay_" + uuid.NewString()

	// Create payment object
	// In CreatePayment method, when creating payment object:
	payment := &models.Payment{
		ID:            paymentID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		Status:        models.PaymentStatusPending,
		PaymentMethod: paymentMethod,
		Description:   sql.NullString{String: description, Valid: description != ""}, // ← Handle NULL
	}

	// Save to database
	if err := s.repo.Create(ctx, payment); err != nil {
		logger.Error(ctx, method, err)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	logger.Info(ctx, method, "Payment created: "+paymentID)

	return payment, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, paymentID string) (*models.Payment, error) {
	ctx, method := logger.FuncInitializer(ctx, "GetPayment")
	defer logger.FuncDisposer(ctx, method)

	logger.Info(ctx, method, "Fetching payment: "+paymentID)

	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		logger.Error(ctx, method, err)
		return nil, paymenterrors.ErrPaymentNotFound
	}

	ctx = logger.WithUserID(ctx, payment.UserID)
	logger.Info(ctx, method, "Payment found")

	return payment, nil
}

// GetUserPayments retrieves all payments for a user
func (s *PaymentService) GetUserPayments(ctx context.Context, userID int, limit, offset int) ([]*models.Payment, int, error) {
	ctx, method := logger.FuncInitializer(ctx, "GetUserPayments")
	defer logger.FuncDisposer(ctx, method)

	ctx = logger.WithUserID(ctx, userID)

	logger.Info(ctx, method, fmt.Sprintf("Fetching payments: limit=%d, offset=%d", limit, offset))

	// Validate pagination
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	payments, total, err := s.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		logger.Error(ctx, method, err)
		return nil, 0, fmt.Errorf("failed to fetch payments: %w", err)
	}

	logger.Info(ctx, method, fmt.Sprintf("Found %d payments (total: %d)", len(payments), total))

	return payments, total, nil
}

// UpdatePaymentStatus updates a payment status
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, paymentID, status string) error {
	ctx, method := logger.FuncInitializer(ctx, "UpdatePaymentStatus")
	defer logger.FuncDisposer(ctx, method)

	logger.Info(ctx, method, fmt.Sprintf("Updating payment %s to status: %s", paymentID, status))

	if err := s.repo.UpdateStatus(ctx, paymentID, status); err != nil {
		logger.Error(ctx, method, err)
		return err
	}

	logger.Info(ctx, method, "Payment status updated")

	return nil
}

// Helper function to validate currency
func isValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		"USD": true,
		"EUR": true,
		"GBP": true,
		"INR": true,
		"JPY": true,
	}
	return validCurrencies[currency]
}
