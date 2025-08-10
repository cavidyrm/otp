package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"otp/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(config *config.Config) (*Database, error) {
	db, err := sql.Open("postgres", config.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			phone_number VARCHAR(20) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			last_login_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS otps (
			id SERIAL PRIMARY KEY,
			phone_number VARCHAR(20) NOT NULL,
			code VARCHAR(10) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			used BOOLEAN DEFAULT FALSE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number)`,
		`CREATE INDEX IF NOT EXISTS idx_otps_phone_number ON otps(phone_number)`,
		`CREATE INDEX IF NOT EXISTS idx_otps_expires_at ON otps(expires_at)`,
		`CREATE INDEX IF NOT EXISTS idx_otps_created_at ON otps(created_at)`,
	}

	for _, query := range queries {
		_, err := d.DB.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}
