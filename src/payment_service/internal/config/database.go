package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase(config *Config) (*sql.DB, error) {
	// build connection string
	fmt.Printf("DB Config: %+v\n", config.Database)
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Name, config.Database.SSLMode)

	fmt.Println("Connection String:", connStr)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connection establised successfully")
	return db, nil
}
