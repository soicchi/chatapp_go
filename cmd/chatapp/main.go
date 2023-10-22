package main

import (
	"log"
	"os"

	"chatapp/internal/infrastructure/database"
)

func main() {
	// Create a new database connection
	dbConfig := database.NewPostgresConfig(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)
	if err := dbConfig.Validate(); err != nil {
		log.Fatal(err)
	}

	dsn := dbConfig.NewDSN()
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database:", db.Name())

	// Migrate the database
	database.Migrate(db)
	log.Println("Successfully migrated database")
}
