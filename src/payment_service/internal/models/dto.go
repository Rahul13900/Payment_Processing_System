package models

// CreatePaymentRequest is the request body for creating a payment
type CreatePaymentRequest struct {
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Description   string `json:"description" binding:"max=500"`
}

// PaymentResponse is returned for payment operations
type PaymentResponse struct {
	ID            string `json:"id"`
	UserID        int    `json:"user_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	CreatedAt     string `json:"created_at"`
}

// PaymentListResponse is returned for listing payments
type PaymentListResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int               `json:"total"`
}

// ErrorResponse for errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
