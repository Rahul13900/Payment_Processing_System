package main

import (
	"fmt"
	"log"
	"os"
	"payment_service/internal/config"
	"payment_service/internal/handler"
	"payment_service/internal/repository"
	"payment_service/internal/service"
	"shared/jwt"
	"shared/logger"
	"shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env.local"); err != nil {
		log.Println("Note: .env.local file not found, using default values")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Bootstrap logging
	ctx := logger.NewRequestContext(nil, "startup")
	logger.Info(ctx, "main", "Payment Service starting")

	// Connect to database
	db, err := config.NewDatabase(cfg)
	if err != nil {
		logger.Error(ctx, "main", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info(ctx, "main", "Database connected")

	// Initialize layers
	paymentRepo := repository.NewPostgresPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.Expiry,
	)

	// Create Gin router
	router := gin.Default()

	// Global middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Public routes
	router.GET("/health", healthHandler)

	// Protected routes (require authentication)
	payments := router.Group("/api/v1/payments")
	payments.Use(middleware.AuthMiddleware(jwtManager))
	{
		payments.POST("", paymentHandler.CreatePayment)
		payments.GET("", paymentHandler.GetUserPayments)
		payments.GET("/:id", paymentHandler.GetPayment)
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info(ctx, "main", fmt.Sprintf("Payment Service listening on %s", addr))

	if err := router.Run(addr); err != nil {
		logger.Error(ctx, "main", err)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "payment-service",
	})
}
