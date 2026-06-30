package models

import (
	"database/sql"
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID                    string         `json:"id" db:"id"`
	UserID                int            `json:"user_id" db:"user_id"`
	Amount                int64          `json:"amount" db:"amount"`
	Currency              string         `json:"currency" db:"currency"`
	Status                string         `json:"status" db:"status"`
	PaymentMethod         string         `json:"payment_method" db:"payment_method"`
	ProviderTransactionID sql.NullString `json:"provider_transaction_id,omitempty" db:"provider_transaction_id"`
	Description           sql.NullString `json:"description,omitempty" db:"description"`
	CreatedAt             time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at" db:"updated_at"`
}

// Payment Status constants
const (
	PaymentStatusPending    = "PENDING"
	PaymentStatusProcessing = "PROCESSING"
	PaymentStatusSucceeded  = "SUCCEEDED"
	PaymentStatusFailed     = "FAILED"
	PaymentStatusRefunded   = "REFUNDED"
	PaymentStatusCancelled  = "CANCELLED"
)

// Payment Method constants
const (
	PaymentMethodCard = "card"
	PaymentMethodUPI  = "upi"
)
