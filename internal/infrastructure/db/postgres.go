package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

// NewPostgresDB creates a new postgres database connection
func NewPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

// newDSN creates a new postgres data source name
func (c *PostgresConfig) NewDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
	)
}

func (c *PostgresConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("postgres host cannot be empty")
	}
	if c.User == "" {
		return fmt.Errorf("postgres user cannot be empty")
	}
	if c.Password == "" {
		return fmt.Errorf("postgres password cannot be empty")
	}
	if c.DBName == "" {
		return fmt.Errorf("postgres dbname cannot be empty")
	}
	if c.Port == "" {
		return fmt.Errorf("postgres port cannot be empty")
	}
	if c.SSLMode == "" {
		return fmt.Errorf("postgres sslmode cannot be empty")
	}

	return nil
}
