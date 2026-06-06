package main

import (
	"fmt"
	"log"
	"user_service/internal/config"
	"user_service/internal/handler"
	"user_service/internal/repository"
	"user_service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env.local"); err != nil {
		log.Println("Note: env file not found, using default values")
	}

	// load config
	cfg := config.LoadConfig()

	// Database connection
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}
	defer db.Close()

	// Initialize layers
	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	UserHandler := handler.NewUserHandler(userService)

	// router
	router := gin.Default()

	// Auth Routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", UserHandler.Register)
		auth.POST("/login", UserHandler.Login)
	}

	// user routes
	users := router.Group("/api/v1/users")
	{
		users.GET("/me", UserHandler.GetProfile)
	}

	// start server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	fmt.Println("User service listening on", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server", err)
	}

}
