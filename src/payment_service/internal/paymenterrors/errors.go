package paymenterrors

import "errors"

var (
	ErrInvalidAmount   = errors.New("invalid amount: must be greater than 0")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidStatus   = errors.New("invalid payment status")
	ErrUserNotFound    = errors.New("user not found")
	ErrInternalError   = errors.New("internal server error")
)
