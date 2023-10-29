package database

import (
	"fmt"

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

func NewPostgresConfig(host, user, password, dbName, port, sslMode string) (*PostgresConfig, error) {
	configFields := map[string]string{
		"host":     host,
		"user":     user,
		"password": password,
		"dbname":   dbName,
		"port":     port,
		"sslmode":  sslMode,
	}
	for fieldName, fieldValue := range configFields {
		if fieldValue == "" {
			return nil, fmt.Errorf("postgres %s cannot be empty", fieldName)
		}
	}

	return &PostgresConfig{
		Host:     host,
		User:     user,
		Password: password,
		DBName:   dbName,
		Port:     port,
		SSLMode:  sslMode,
	}, nil
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

func Migrate(db *gorm.DB, entities ...interface{}) error {
	return db.AutoMigrate(entities...)
}
