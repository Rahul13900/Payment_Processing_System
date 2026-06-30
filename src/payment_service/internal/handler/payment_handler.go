package handler

import (
	"net/http"
	"payment_service/internal/models"
	"payment_service/internal/paymenterrors"
	"payment_service/internal/service"
	"shared/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

// CreatePayment handles POST /payments
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	ctx, method := logger.FuncInitializer(c.Request.Context(), "CreatePayment")
	defer logger.FuncDisposer(ctx, method)

	var req models.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, method, "Invalid request body")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		logger.Warn(ctx, method, "User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		logger.Warn(ctx, method, "Invalid user ID type")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Invalid user ID",
		})
		return
	}

	ctx = logger.WithUserID(ctx, userID)

	// Call service
	payment, err := h.service.CreatePayment(
		c.Request.Context(),
		userID,
		req.Amount,
		req.Currency,
		req.PaymentMethod,
		req.Description,
	)

	if err != nil {
		logger.Error(ctx, method, err)

		if err == paymenterrors.ErrInvalidAmount {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "INVALID_AMOUNT",
				Message: "Amount must be greater than 0",
			})
			return
		}

		if err == paymenterrors.ErrInvalidCurrency {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "INVALID_CURRENCY",
				Message: "Invalid currency",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to create payment",
		})
		return
	}

	resp := models.PaymentResponse{
		ID:            payment.ID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt.String(),
	}

	c.JSON(http.StatusCreated, resp)
}

// GetPayment handler
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	ctx, method := logger.FuncInitializer(c.Request.Context(), "GetPayment")
	defer logger.FuncDisposer(ctx, method)

	paymentID := c.Param("id")

	payment, err := h.service.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		logger.Error(ctx, method, err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Payment not found",
		})
		return
	}

	ctx = logger.WithUserID(ctx, payment.UserID)

	resp := models.PaymentResponse{
		ID:            payment.ID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt.String(),
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserPayments handler
func (h *PaymentHandler) GetUserPayments(c *gin.Context) {
	ctx, method := logger.FuncInitializer(c.Request.Context(), "GetUserPayments")
	defer logger.FuncDisposer(ctx, method)

	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		logger.Warn(ctx, method, "User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Invalid user ID",
		})
		return
	}

	ctx = logger.WithUserID(ctx, userID)

	// Parse pagination params
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	payments, total, err := h.service.GetUserPayments(c.Request.Context(), userID, limit, offset)
	if err != nil {
		logger.Error(ctx, method, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to fetch payments",
		})
		return
	}

	var respPayments []models.PaymentResponse
	for _, p := range payments {
		respPayments = append(respPayments, models.PaymentResponse{
			ID:            p.ID,
			UserID:        p.UserID,
			Amount:        p.Amount,
			Currency:      p.Currency,
			Status:        p.Status,
			PaymentMethod: p.PaymentMethod,
			CreatedAt:     p.CreatedAt.String(),
		})
	}

	c.JSON(http.StatusOK, models.PaymentListResponse{
		Payments: respPayments,
		Total:    total,
	})
}
