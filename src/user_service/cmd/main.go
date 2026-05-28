package main

import (
	"fmt"
	"log"
	"user_service/internal/config"

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

	// router
	router := gin.Default()

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	fmt.Println("User service listening on", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server", err)
	}

}
